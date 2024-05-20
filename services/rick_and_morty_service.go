package services

import (
	"VSK_test/api"
	"VSK_test/types"
	"VSK_test/utils"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"sync"
)

const CharacterPath = "/character"
const EpisodePath = "/episode"

type APIService interface {
	GetCharacter(id int, ch chan types.Character, wg *sync.WaitGroup)
	GetEpisode(id int, ch chan types.Episode, wg *sync.WaitGroup)
	Run()
}

type RickAndMortyService struct {
	client *api.HttpClient
}

// NewRickAndMortyService creates a new instance of the RickAndMortyService struct.
//
// It takes a pointer to an HttpClient from the api package as a parameter and returns an APIService interface.
func NewRickAndMortyService(HttpClient *api.HttpClient) APIService {
	return &RickAndMortyService{
		client: HttpClient,
	}
}

// GetCharacter fetches a character by ID from the RickAndMorty API and sends it through the provided channel.
//
// Parameters:
// - id: the ID of the character to fetch.
// - ch: the channel to send the fetched character through.
// - wg: the wait group to signal when the fetching is done.
func (k *RickAndMortyService) GetCharacter(id int, ch chan types.Character, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := k.client.Get(fmt.Sprintf("%s/%d", CharacterPath, id))
	if err != nil {
		log.Printf("Error fetching character %d: %v\n", id, err)
		return
	}

	character := new(types.Character)
	if err := mapstructure.Decode(resp, &character); err != nil {
		log.Printf("Error decoding character %d: %v\n", id, err)
		return
	}

	ch <- *character
}

// GetEpisode fetches an episode by ID from the RickAndMorty API and sends it through the provided channel.
//
// Parameters:
// - id: the ID of the episode to fetch.
// - ch: the channel to send the fetched episode through.
// - wg: the wait group to signal when the fetching is done.
func (k *RickAndMortyService) GetEpisode(id int, ch chan types.Episode, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := k.client.Get(fmt.Sprintf("%s/%d", EpisodePath, id))
	if err != nil {
		log.Printf("Error fetching episode %d: %v\n", id, err)
		return
	}

	episode := new(types.Episode)

	if err := mapstructure.Decode(res, &episode); err != nil {
		log.Printf("Error decoding episode %d: %v\n", id, err)
		return
	}
	ch <- *episode
}

// Run executes the main logic of the RickAndMortyService.
//
// It fetches episodes from the RickAndMorty API and retrieves the second-to-last character of each episode.
// The episodes and their corresponding characters are stored in channels and processed concurrently.
// The function then prints the JSON representation of each episode with the second-to-last character.
func (k *RickAndMortyService) Run() {
	var wg sync.WaitGroup
	episodesChan := make(chan types.Episode, 10)
	characterChan := make(chan types.Character, 10)
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go k.GetEpisode(i, episodesChan, &wg)
	}
	wg.Wait()
	close(episodesChan)

	var characterWg sync.WaitGroup
	var episodes []types.Episode

	for episode := range episodesChan {
		if len(episode.Characters) < 2 {
			log.Printf("Episode %d: Not enough characters\n", episode.ID)
			continue
		}
		episodes = append(episodes, episode)
		characterWg.Add(1)
		characterId, err := utils.GetIdFromUrl(episode.Characters[len(episode.Characters)-2])
		if err != nil {
			log.Printf("Error converting ID to int: %v\n", err)
		}
		go k.GetCharacter(characterId, characterChan, &characterWg)
	}

	characterWg.Wait()
	close(characterChan)
	characterMap := make(map[int]string)
	for character := range characterChan {
		characterMap[character.ID] = character.Name
	}
	for i, episode := range episodes {
		if len(episode.Characters) < 2 {
			continue
		}
		characterId, err := utils.GetIdFromUrl(episode.Characters[len(episode.Characters)-2])
		if err != nil {
			log.Printf("Error converting ID to int: %v\n", err)
		}
		episodes[i].SecondLastCharacter = characterMap[characterId]
		episodeJSON, err := json.MarshalIndent(episodes[i], "", "  ")
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(string(episodeJSON))
		fmt.Println("\n#################################################################\n")
	}

}
