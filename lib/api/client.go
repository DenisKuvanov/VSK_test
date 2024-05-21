package api

import (
	"VSK_test/lib"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HttpClient struct {
	engine  *http.Client
	baseURL string
}

// NewHttpClient creates a new instance of HttpClient with the provided ClientSettings.
//
// Parameters:
// - config: the ClientSettings to configure the HttpClient.
//
// Returns:
// - *HttpClient: a pointer to the newly created HttpClient instance.
func NewHttpClient(config lib.ClientSettings) *HttpClient {
	return &HttpClient{
		engine: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		baseURL: config.BaseUrl,
	}
}

// Get sends a GET request to the specified path and returns the response body as an interface{}.
//
// Parameters:
// - path: the path to send the GET request to.
//
// Returns:
// - interface{}: the response body as an interface{}.
// - error: an error if the request fails.
func (c *HttpClient) Get(path string) (interface{}, error) {
	res, err := c.engine.Get(fmt.Sprintf("%s%s", c.baseURL, path))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	var response interface{}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
