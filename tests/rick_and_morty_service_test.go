package tests

import (
	"VSK_test/api"
	"VSK_test/services"
	"VSK_test/types"
	"sync"
	"testing"
)

func TestRickAndMortyService_GetCharacter(t *testing.T) {
	client := api.NewHttpClient()
	service := services.NewRickAndMortyService(client)

	ch := make(chan types.Character, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go service.GetCharacter(1, ch, &wg)
	wg.Wait()
	close(ch)

	character := <-ch
	if character.ID != 1 {
		t.Errorf("Expected ID to be 1, got %v", character.ID)
	}
}

func TestRickAndMortyService_GetEpisode(t *testing.T) {
	client := api.NewHttpClient()
	service := services.NewRickAndMortyService(client)

	ch := make(chan types.Episode, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go service.GetEpisode(1, ch, &wg)
	wg.Wait()
	close(ch)

	episode := <-ch
	if episode.ID != 1 {
		t.Errorf("Expected ID to be 1, got %v", episode.ID)
	}
}
