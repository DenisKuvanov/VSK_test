package api

import (
	"VSK_test/lib"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHttpClient_Get tests the Get method of HttpClient
func TestHttpClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"key": "value"}`))
	}))
	defer server.Close()

	client := NewHttpClient(lib.ClientSettings{
		Timeout: 5,
		BaseUrl: server.URL,
	})

	// Test successful GET request
	resp, err := client.Get("/test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.(map[string]interface{})["key"] != "value" {
		t.Errorf("Expected response 'value', got %v", resp.(map[string]interface{})["key"])
	}

	// Test 404 error
	_, err = client.Get("/bad-test")
	if err == nil {
		t.Fatalf("Expected error for bad-test path, got none")
	}
	if err.Error() != "status code error: 404 404 Not Found" {
		t.Errorf("Expected 404 error, got %v", err)
	}
}
