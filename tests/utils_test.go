package tests

import (
	"VSK_test/utils"
	"testing"
)

func TestGetIdFromUrl(t *testing.T) {
	tests := []struct {
		url      string
		expected int
	}{
		{"https://rickandmortyapi.com/api/character/1", 1},
		{"https://rickandmortyapi.com/api/character/42", 42},
		{"https://rickandmortyapi.com/api/episode/3", 3},
	}

	for _, test := range tests {
		result, err := utils.GetIdFromUrl(test.url)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != test.expected {
			t.Errorf("For URL %s, expected %d, but got %d", test.url, test.expected, result)
		}
	}
}
