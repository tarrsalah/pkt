package pkt

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetRequestToken(t *testing.T) {
	// request
	consumerKey := "123"
	redirectUri := "http://localhost"

	// response
	expectedCode := "456"

	client, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/oauth/request", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		code := struct {
			Code string `json:"code"`
		}{Code: expectedCode}
		json.NewEncoder(w).Encode(code)
	})

	code, err := client.getRequestToken(consumerKey, redirectUri)
	if err != nil {
		t.Error(err)
	}

	if code != expectedCode {
		t.Errorf("Expected code to be %s, got %s", expectedCode, code)
	}
}

func TestGetAccessToken(t *testing.T) {
	// request
	consumerKey := "123"
	code := "456"

	// response
	expectedAccessToken := "789"
	expectedUserName := "user1"

	client, mux, server := testServer()
	defer server.Close()

	mux.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		accessToken := struct {
			AccessToken string `json:"access_token"`
			Username    string `json:"username"`
		}{
			AccessToken: expectedAccessToken,
			Username:    expectedUserName,
		}
		json.NewEncoder(w).Encode(accessToken)
	})

	accessToken, err := client.getAccessToken(consumerKey, code)
	if err != nil {
		t.Error(err)
	}

	if accessToken != expectedAccessToken {
		t.Errorf("Expected access token to be %s, but got %s", expectedAccessToken, accessToken)
	}

}
