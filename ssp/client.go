package ssp

import (
	"encoding/json"
	"fmt"
	"github.com/google/jsonapi"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	Config  *Config
	baseURL *url.URL
	client  *http.Client
}

type ErrorResponse struct {
	Errors []struct {
		Status string `json:"status"`
		Title  string `json:"title"`
	} `json:"errors"`
}

func (er *ErrorResponse) String() string {
	messages := make([]string, len(er.Errors))
	for i, e := range er.Errors {
		messages[i] = e.Title
	}
	return strings.Join(messages, ", ")
}

func NewClient(c *Config) (*Client, error) {
	if c == nil {
		c = NewDefaultConfig()
	}

	parsed, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}

	a := &Client{
		Config:  c,
		baseURL: parsed,
		client:  c.Client(),
	}

	return a, nil
}

func (a *Client) get(path string) (io.ReadCloser, error) {
	resp, err := a.request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (a *Client) post(path string, body io.Reader) (io.ReadCloser, error) {
	resp, err := a.request("POST", path, body)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (a *Client) delete(path string, body io.Reader) (io.ReadCloser, error) {
	resp, err := a.request("DELETE", path, body)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (a *Client) request(method string, path string, body io.Reader) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", a.baseURL, path)
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", jsonapi.MediaType)
	req.Header.Add("Accept", jsonapi.MediaType)
	req.Header.Add("X-Api-Version", "2.0")

	if os.Getenv("DEBUG") != "" {
		dump, _ := httputil.DumpRequestOut(req, true)
		fmt.Printf("%s", dump)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		er := &ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(er); err == nil {
			return nil, fmt.Errorf("HTTP %d - '%s'", resp.StatusCode, er)
		} else {
			return nil, fmt.Errorf("HTTP %d - '%s'", resp.StatusCode, resp.Status)
		}
	}

	if os.Getenv("DEBUG") != "" {
		dump, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("%s", dump)
	}

	if resp.Header.Get("Content-Type") != "application/vnd.api+json" {
		return nil, fmt.Errorf("Unexpected Content-Type: '%s'", resp.Header.Get("Content-Type"))
	}

	return resp, nil
}

func parseSSTime(t string) (time.Time, error) {
	formats := []string{
		"15:04",
		"15:04:05",
	}

	var err error
	var parsed time.Time
	for _, format := range formats {
		parsed, err = time.Parse(format, t)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, err
}