package api

import (
	"VSK_test/lib"
	"testing"
)

func TestHttpClient_GetSuccess(t *testing.T) {
	config := lib.LoadConfig()
	client := NewHttpClient(config.ClientConfig)
	_, err := client.Get("/character/1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHttpClient_GetError(t *testing.T) {
	config := lib.LoadConfig()
	client := NewHttpClient(config.ClientConfig)
	_, err := client.Get("/character/0")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
