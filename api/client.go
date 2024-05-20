package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://rickandmortyapi.com/api"

type HttpClient struct {
	engine *http.Client
}

// NewHttpClient creates a new instance of HttpClient with a timeout of 3 seconds.
//
// It returns a pointer to the newly created HttpClient.
func NewHttpClient() *HttpClient {
	return &HttpClient{
		engine: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// Get retrieves data from the specified path using the HttpClient.
//
// It takes a string parameter `path` which represents the path to be appended to the baseURL.
// The function returns an `interface{}` and an `error`. The `interface{}` represents the response data,
// and the `error` represents any error that occurred during the request.
func (c *HttpClient) Get(path string) (interface{}, error) {
	res, err := c.engine.Get(fmt.Sprintf("%s%s", baseURL, path))
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
