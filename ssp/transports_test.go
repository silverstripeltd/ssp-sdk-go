package ssp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type basicAuthReflector struct{}

func (basicAuthReflector) RoundTrip(req *http.Request) (*http.Response, error) {
	username, password, _ := req.BasicAuth()
	res := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(strings.NewReader(fmt.Sprintf("%s:%s", username, password))),
	}
	return res, nil
}

func TestRoundTrip(t *testing.T) {
	r := &basicAuthReflector{}
	bat := &BasicAuthTransport{
		Username:  "testUsername",
		Password:  "testPassword",
		Transport: r,
	}

	res, err := bat.RoundTrip(&http.Request{})
	if err != nil {
		t.Errorf("RoundTrip failed: %s", err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	expected := fmt.Sprintf("%s:%s", "testUsername", "testPassword")
	if string(body) != expected {
		t.Errorf("Username or password were not properly added to the request ('%s' vs '%s')", body, expected)

	}
}
