package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/pkg/browser"
)

const baseUrl = "https://getpocket.com/v3"

type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       *Creds
}

func NewClient(auth *Creds) *Client {
	return &Client{
		baseURL:    baseUrl,
		httpClient: http.DefaultClient,
		auth:       auth,
	}
}

func (c *Client) sendPost(action string, in interface{}, out interface{}) error {
	url := c.baseURL + action

	requestBody, err := json.Marshal(in)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Accept", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return errors.New(response.Status)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &out)
	// Handle empty response
	if err != nil {
		empty := &struct {
			Status int `json:"status"`
		}{Status: 0}
		json.Unmarshal(responseBody, empty)
		if empty.Status == 2 {
			return nil
		}
	}

	return err
}

// Creds is a pair of keys
type Creds struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}

// Authenticate a consumer key
func (c *Client) Authenticate(consumerKey string) *Creds {
	redirectURI := fmt.Sprintf("http://localhost:%d/oauth/pocket/callback", 5000)
	requestToken, err := c.getRequestToken(consumerKey, redirectURI)
	if err != nil {
		log.Fatal(err)
	}

	authorizeURL := fmt.Sprintf(
		"https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
		requestToken,
		redirectURI,
	)

	ch := make(chan struct{})

	server := http.Server{
		Addr: fmt.Sprintf(":%v", 5000),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Authentication succeded!")
			ch <- struct{}{}
		}),
	}

	go func() {
		server.ListenAndServe()
	}()
	defer server.Close()

	browser.OpenURL(authorizeURL)
	<-ch

	accessToken, err := c.getAccessToken(consumerKey, requestToken)
	if err != nil {
		log.Fatal(err)
	}

	// Authentication succeeded
	auth := &Creds{}
	auth.ConsumerKey = consumerKey
	auth.AccessToken = accessToken
	return auth
}

func (c *Client) getRequestToken(consumerKey, redirectURI string) (string, error) {
	in := struct {
		ConsumerKey string `json:"consumer_key"`
		RedirectURI string `json:"redirect_uri"`
	}{
		ConsumerKey: consumerKey,
		RedirectURI: redirectURI,
	}

	var out struct {
		Code string `json:"code"`
	}

	err := c.sendPost("/oauth/request", in, &out)
	return out.Code, err
}

func (c *Client) getAccessToken(consumerKey, code string) (string, error) {
	in := struct {
		ConsumerKey string `json:"consumer_key"`
		Code        string `json:"code"`
	}{
		ConsumerKey: consumerKey,
		Code:        code,
	}

	var out struct {
		AccessToken string `json:"access_token"`
		Username    string `json:"username"`
	}

	err := c.sendPost("/oauth/authorize", in, &out)
	return out.AccessToken, err
}

func (c *Client) Retrieve(after string, offset int) ([]Item, error) {
	action := "/get"
	request := retrieveRequest{
		Creds: *c.auth,
		retrieveOptions: retrieveOptions{
			Since:      after,
			Offset:     offset,
			DetailType: "complete",
			Sort:       "newest",
			Count:      PageCount,
		},
	}
	response := retrieveResponse{}

	err := c.sendPost(action, request, &response)

	return response.Items(), err
}

func (c *Client) RetrieveAll(after string) ([]Item, error) {
	items := []Item{}
	offset := 0
	for {
		retrieved, err := c.Retrieve(after, offset)
		if err != nil {
			return nil, err
		}

		if len(retrieved) == 0 {
			break
		}

		items = append(items, retrieved...)
		offset = offset + PageCount
	}

	return items, nil
}

type retrieveOptions struct {
	DetailType string `json:"detailType"`
	Since      string `json:"since"`
	Offset     int    `json:"offset"`
	Sort       string `json:"sort"`
	Count      int    `json:"count"`
}

type retrieveRequest struct {
	Creds
	retrieveOptions
}

type retrieveResponse struct {
	Status int             `json:"status"`
	List   map[string]Item `json:"list"`
}

func (r retrieveResponse) Items() []Item {
	var items []Item

	for _, item := range r.List {
		items = append(items, item)
	}

	sort.Slice(items[:], func(i, j int) bool {
		return items[i].AddedAt >= items[j].AddedAt
	})

	return items
}
