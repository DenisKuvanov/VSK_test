package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	rnm "github.com/pitakill/rickandmortyapigowrapper"
)

type MyEpisode struct {
	rnm.Episode
	SecondLastCharacter string
}

func fetchEpisode(id int, ch chan<- rnm.Episode, wg *sync.WaitGroup) {
	defer wg.Done()
	episode, err := rnm.GetEpisode(id)
	if err != nil {
		log.Printf("Error fetching episode %d: %v\n", id, err)
		return
	}
	ch <- *episode
}

func fetchCharacter(url string, ch chan<- rnm.Character) {
	characterIdStr := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	characterId, err := strconv.Atoi(characterIdStr)
	if err != nil {
		log.Printf("Error converting character ID %s to int: %v\n", characterIdStr, err)
		return
	}
	character, err := rnm.GetCharacter(characterId)
	if err != nil {
		log.Printf("Error fetching character %d: %v\n", characterId, err)
		return
	}
	ch <- *character
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	episodesChan := make(chan rnm.Episode, 10)
	characterChan := make(chan rnm.Character, 10)

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go fetchEpisode(i, episodesChan, &wg)
	}
	wg.Wait()
	close(episodesChan)

	var characterWg sync.WaitGroup
	var episodes []MyEpisode

	for episode := range episodesChan {
		if len(episode.Characters) < 2 {
			log.Printf("Episode %d: Not enough characters\n", episode.ID)
			continue
		}
		episodes = append(episodes, MyEpisode{Episode: episode})

		characterWg.Add(1)
		go func(characterURL string) {
			defer characterWg.Done()
			fetchCharacter(characterURL, characterChan)
		}(episode.Characters[len(episode.Characters)-2])
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
		characterURL := episode.Characters[len(episode.Characters)-2]
		characterIdStr := strings.Split(characterURL, "/")[len(strings.Split(characterURL, "/"))-1]
		characterId, _ := strconv.Atoi(characterIdStr)
		episodes[i].SecondLastCharacter = characterMap[characterId]

		episodeJSON, err := json.MarshalIndent(episodes[i], "", "  ")
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(string(episodeJSON))
		fmt.Println("\n#################################################################\n")
	}

	duration := time.Since(start)
	fmt.Printf("Script took %s\n", duration)
}
