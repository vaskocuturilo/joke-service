package main

import (
	"context"
	"fmt"
	"joke-api-service/internal/api"
	"joke-api-service/internal/config"
	"joke-api-service/internal/domain"
	"log"
	"time"
)

func main() {

	var timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	appConfig := config.Config{BaseUrl: "https://official-joke-api.appspot.com"}

	var provider domain.JokeProvider = api.NewClient(appConfig.BaseUrl, timeout)

	result, err := provider.GetTenJokes(ctx)

	if err != nil {
		log.Fatalf("Critical error fetching joke: %v", err)
	}

	for count, value := range result {
		fmt.Printf("%d. %s, %s \n", count+1, value.Setup, value.Punchline)
	}
}
