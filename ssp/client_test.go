package ssp

import (
	"fmt"
	"github.com/google/jsonapi"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func ExampleClient() {
	// Default client uses either the $HOME/.dashboard.env, or environment variable overrides.
	ssp, err := NewClient(nil)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	env, _ := ssp.GetEnvironment("mystack", "myenv")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	fmt.Printf("Just fetched environment %s", env.Name)
}

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
