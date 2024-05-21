package services

import (
	"VSK_test/lib"
	"VSK_test/lib/api"
	"VSK_test/lib/types"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// createMockServer creates a mock server with predefined responses
func createMockServer(t *testing.T, paths map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, ok := paths[r.URL.Path]
		if !ok {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
}

func TestRickAndMortyService_GetCharacter(t *testing.T) {
	mockResponses := map[string]string{
		"/character/1": `{"id": 1, "name": "Rick Sanchez"}`,
	}

	server := createMockServer(t, mockResponses)
	defer server.Close()

	client := api.NewHttpClient(lib.ClientSettings{
		Timeout: 5,
		BaseUrl: server.URL,
	})

	service := NewRickAndMortyService(client, lib.ServiceSettings{
		CharacterPath: "/character",
		EpisodePath:   "/episode",
	})

	rms := service.(*RickAndMortyService)
	rms.characterChan = make(chan types.Character, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	rms.characterWg.Add(1)
	go func() {
		defer wg.Done()
		rms.GetCharacter(1)
	}()
	wg.Wait()

	select {
	case character := <-rms.characterChan:
		if character.ID != 1 || character.Name != "Rick Sanchez" {
			t.Errorf("Expected character Rick Sanchez, got %v", character)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Character not received on channel")
	}
}

func TestRickAndMortyService_GetEpisode(t *testing.T) {
	mockResponses := map[string]string{
		"/episode/1":   `{"id": 1, "name": "Pilot", "characters": ["/character/1", "/character/2", "/character/3"]}`,
		"/character/2": `{"id": 2, "name": "Morty Smith"}`,
	}

	server := createMockServer(t, mockResponses)
	defer server.Close()

	client := api.NewHttpClient(lib.ClientSettings{
		Timeout: 5,
		BaseUrl: server.URL,
	})

	service := NewRickAndMortyService(client, lib.ServiceSettings{
		CharacterPath: "/character",
		EpisodePath:   "/episode",
	})

	rms := service.(*RickAndMortyService)
	rms.episodeChan = make(chan types.Episode, 1)
	rms.characterChan = make(chan types.Character, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	rms.episodeWg.Add(1)
	go func() {
		defer wg.Done()
		rms.GetEpisode(1)
	}()
	wg.Wait()

	select {
	case episode := <-rms.episodeChan:
		if episode.ID != 1 || episode.Name != "Pilot" {
			t.Errorf("Expected episode Pilot, got %v", episode)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Episode not received on channel")
	}

	select {
	case character := <-rms.characterChan:
		if character.ID != 2 || character.Name != "Morty Smith" {
			t.Errorf("Expected character Morty Smith, got %v", character)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Character not received on channel")
	}
}
