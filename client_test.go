package pkt

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

	client.post("/", in, &out)
	if out.Id != 1 {
		t.Fatalf("response body id expected to be 1, got %d", out.Id)
	}
}

func testServer() (*Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := NewClient(nil)
	client.baseURL = server.URL

	return client, mux, server
}
