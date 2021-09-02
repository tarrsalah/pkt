package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer() (*Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient(nil)
	client.baseURL = server.URL

	return client, mux, server
}

func TestClient(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()

	type intype struct {
		Id int
	}

	type outtype struct {
		Id int
	}

	mux.HandleFunc(("/"), func(w http.ResponseWriter, r *http.Request) {
		var in intype
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if in.Id != 0 {
			t.Fatalf("request body id expected to be 0, got %d", in.Id)
		}

		out := outtype{
			Id: 1,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)

	})

	in := intype{Id: 0}
	var out outtype

	client.sendPost("/", in, &out)
	if out.Id != 1 {
		t.Fatalf("response body id expected to be 1, got %d", out.Id)
	}
}

func TestGetRequestToken(t *testing.T) {
	// request
	consumerKey := "123"
	redirectURI := "http://localhost"

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

	code, err := client.getRequestToken(consumerKey, redirectURI)
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
