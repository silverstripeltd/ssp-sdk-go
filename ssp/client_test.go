package ssp

import (
	"github.com/google/jsonapi"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewApi(t *testing.T) {

	api, _ := NewClient(&Config{
		BaseURL: "http://localhost",
	})

	if api.baseURL.Host != "localhost" {
		t.Error("BaseURL not parsed correctly")
	}
}

func newMockDashboard(in interface{}, responseCode int) (*Client, *httptest.Server) {

	responseHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", jsonapi.MediaType)
		w.WriteHeader(responseCode)
		jsonapi.MarshalPayload(w, in)
	}

	ts := httptest.NewServer(http.HandlerFunc(responseHandler))

	api, err := NewClient(&Config{BaseURL: ts.URL})
	if err != nil {
		panic("Something bad happened while creating the API")
	}
	return api, ts
}
