package main

import (
	"VSK_test/api"
	"VSK_test/services"
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	httpClient := api.NewHttpClient()
	RickAndMortyService := services.NewRickAndMortyService(httpClient)
	RickAndMortyService.Run()

	duration := time.Since(start)
	fmt.Printf("Script took %s\n", duration)
}
