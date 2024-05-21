package main

import (
	"VSK_test/lib"
	"VSK_test/lib/api"
	"VSK_test/lib/log"
	"VSK_test/lib/services"
	"time"
)

func main() {
	start := time.Now()

	config := lib.LoadConfig()
	httpClient := api.NewHttpClient(config.ClientConfig)
	RickAndMortyService := services.NewRickAndMortyService(httpClient, config.ServiceConfig)
	RickAndMortyService.Run()

	duration := time.Since(start)
	log.Log.Infof("Script took %s\n", duration)
}
