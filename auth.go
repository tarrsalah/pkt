package pkt

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/browser"
)

type Auth struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}

type requestTokenRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectUri string `json:"redirect_uri"`
}

type requestTokenResponse struct {
	Code string `json:"code"`
}

type accessTokenRequest struct {
	ConsumerKey string `json:"consumer_key"`
	Code        string `json:"code"`
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}

func (c *Client) Authenticate(consumerKey string) *Auth {
	redirectUri := fmt.Sprintf("http://localhost:%d/oauth/pocket/callback", 5000)
	requestToken, err := c.getRequestToken(consumerKey, redirectUri)
	if err != nil {
		log.Fatal(err)
	}

	authorizeUrl := fmt.Sprintf(
		"https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
		requestToken,
		redirectUri,
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

	browser.OpenURL(authorizeUrl)
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

func (c *Client) getRequestToken(consumerKey, redirectUri string) (string, error) {
	in := requestTokenRequest{
		ConsumerKey: consumerKey,
		RedirectUri: redirectUri,
	}
	var out requestTokenResponse

	err := c.Post("/oauth/request", in, &out)
	return out.Code, err
}

func (c *Client) getAccessToken(consumerKey, code string) (string, error) {
	in := accessTokenRequest{
		ConsumerKey: consumerKey,
		Code:        code,
	}

	var out accessTokenResponse

	err := c.Post("/oauth/authorize", in, &out)
	return out.AccessToken, err
}
