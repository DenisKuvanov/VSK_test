package tests

import (
	"VSK_test/api"
	"testing"
)

func TestHttpClient_GetSuccess(t *testing.T) {
	client := api.NewHttpClient()
	_, err := client.Get("/character/1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHttpClient_GetError(t *testing.T) {
	client := api.NewHttpClient()
	_, err := client.Get("/character/0")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
