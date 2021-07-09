package pkt

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/browser"
)

// Auth is a pair of keys
type Auth struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}

// Authenticate a consumer key
func (c *Client) Authenticate(consumerKey string) *Auth {
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
	auth := &Auth{}
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

	err := c.Post("/oauth/request", in, &out)
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

	err := c.Post("/oauth/authorize", in, &out)
	return out.AccessToken, err
}
