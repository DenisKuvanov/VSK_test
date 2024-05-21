package main

import (
	"VSK_test/lib"
	"VSK_test/lib/api"
	"VSK_test/lib/services"
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	config := lib.LoadConfig()
	httpClient := api.NewHttpClient(config.ClientConfig)
	RickAndMortyService := services.NewRickAndMortyService(httpClient, config.ServiceConfig)
	RickAndMortyService.Run()

	duration := time.Since(start)
	fmt.Printf("Script took %s\n", duration)
}
