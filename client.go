package pkt

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const baseUrl = "https://getpocket.com/v3"

type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       *Auth
}

func NewClient(auth *Auth) *Client {
	return &Client{
		baseURL:    baseUrl,
		httpClient: http.DefaultClient,
		auth:       auth,
	}
}

func (c *Client) Post(action string, in interface{}, out interface{}) error {
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

	res, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return errors.New(res.Status)
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
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
