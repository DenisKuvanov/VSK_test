package services

import (
	"VSK_test/lib"
	"VSK_test/lib/api"
	"VSK_test/lib/log"
	"VSK_test/lib/types"
	"VSK_test/lib/utils"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"sync"
	"time"
)

type APIService interface {
	GetCharacter(id int)
	GetEpisode(id int)
	Run()
}

type RickAndMortyService struct {
	client        *api.HttpClient
	episodeWg     sync.WaitGroup
	characterWg   sync.WaitGroup
	episodeChan   chan types.Episode
	characterChan chan types.Character
	Episodes      []types.Episode
	Characters    map[int]types.Character
	mu            sync.Mutex
	settings      lib.ServiceSettings
}

// NewRickAndMortyService creates a new instance of the RickAndMortyService struct.
//
// Parameters:
// - HttpClient: A pointer to the HttpClient struct.
// - config: An instance of the ServiceSettings struct.
//
// Returns:
// - An implementation of the APIService interface.
func NewRickAndMortyService(HttpClient *api.HttpClient, config lib.ServiceSettings) APIService {
	return &RickAndMortyService{
		client:        HttpClient,
		episodeChan:   make(chan types.Episode),
		characterChan: make(chan types.Character),
		Episodes:      make([]types.Episode, 0),
		Characters:    make(map[int]types.Character),
		settings:      config,
	}
}

// GetCharacter fetches a character with the given ID from the RickAndMorty API and sends it to the character channel.
//
// Parameters:
// - id: The ID of the character to fetch.
func (k *RickAndMortyService) GetCharacter(id int) {
	defer k.characterWg.Done()
	resp, err := k.client.Get(fmt.Sprintf("%s/%d", k.settings.CharacterPath, id))
	if err != nil {
		log.Log.Errorf("Error fetching character %d: %v\n", id, err)
		return
	}

	character := new(types.Character)
	if err := mapstructure.Decode(resp, &character); err != nil {
		log.Log.Errorf("Error decoding character %d: %v\n", id, err)
		return
	}
	k.characterChan <- *character
}

// GetEpisode fetches an episode with the given ID from the RickAndMorty API and sends it to the episode channel.
//
// Parameters:
// - id: The ID of the episode to fetch.
func (k *RickAndMortyService) GetEpisode(id int) {
	defer k.episodeWg.Done()
	res, err := k.client.Get(fmt.Sprintf("%s/%d", k.settings.EpisodePath, id))
	if err != nil {
		log.Log.Errorf("Error fetching episode %d: %v\n", id, err)
		return
	}

	episode := new(types.Episode)

	if err := mapstructure.Decode(res, &episode); err != nil {
		log.Log.Errorf("Error decoding episode %d: %v\n", id, err)
		return
	}
	characterId, err := utils.GetIdFromUrl(episode.Characters[len(episode.Characters)-2])
	if err != nil {
		log.Log.Errorf("Error converting ID to int: %v\n", err)
		return
	}
	k.characterWg.Add(1)
	go k.GetCharacter(characterId)
	k.episodeChan <- *episode
}

// ReadEpisodesChan reads episodes from the episode channel and appends them to the Episodes slice.
//
// This function runs in a loop and continuously listens for episodes on the episode channel. If an episode is received,
// it acquires a lock on the mutex, appends the episode to the Episodes slice, and releases the lock. If the episode
// channel is closed, the function returns. If no episode is received within 5 seconds, the function also returns.
func (k *RickAndMortyService) ReadEpisodesChan() {
	for {
		select {
		case episode, open := <-k.episodeChan:
			if open {
				k.mu.Lock()
				k.Episodes = append(k.Episodes, episode)
				k.mu.Unlock()
			} else {
				return
			}
		case <-time.After(5 * time.Second):
			return
		}
	}
}

// ReadCharacterChan reads characters from the character channel and stores them in the Characters map.
//
// This function runs in a loop and continuously listens for characters on the character channel. If a character is
// received, it acquires a lock on the mutex, stores the character in the Characters map using the character's ID as
// the key, and releases the lock. If the character channel is closed, the function returns. If no character is
// received within 5 seconds, the function also returns.
func (k *RickAndMortyService) ReadCharacterChan() {
	for {
		select {
		case character, open := <-k.characterChan:
			if open {
				k.mu.Lock()
				k.Characters[character.ID] = character
				k.mu.Unlock()
			} else {
				return
			}
		case <-time.After(5 * time.Second):
			return
		}
	}
}

// Run executes the main logic of the RickAndMortyService. It starts two goroutines to read episodes and characters
// from channels, fetches episodes from the API, and populates the Episodes slice. It then iterates over the
// Episodes slice and retrieves the name of the second-to-last character associated with each episode.
// The episode information is printed in JSON format and the length of the Episodes slice is printed.
func (k *RickAndMortyService) Run() {
	go k.ReadEpisodesChan()
	go k.ReadCharacterChan()

	for i := 1; i <= k.settings.EpisodesAmount; i++ {
		k.episodeWg.Add(1)
		go k.GetEpisode(i)
	}
	k.episodeWg.Wait()
	close(k.episodeChan)
	k.characterWg.Wait()
	close(k.characterChan)

	for i, episode := range k.Episodes {
		if len(episode.Characters) < 2 {
			continue
		}
		characterId, err := utils.GetIdFromUrl(episode.Characters[len(episode.Characters)-2])
		if err != nil {
			log.Log.Errorf("Error converting ID to int: %v\n", err)
			continue
		}
		k.Episodes[i].SecondLastCharacter = k.Characters[characterId].Name
		episodeJSON, err := json.MarshalIndent(k.Episodes[i], "", "  ")
		if err != nil {
			log.Log.Errorf("Error marshalling episode: %v\n", err)
			continue
		}
		fmt.Println(string(episodeJSON))
		fmt.Println("\n#################################################################\n")
	}
	fmt.Println("Done!")
	fmt.Printf("Length of episodes: %d\n", len(k.Episodes))
}
