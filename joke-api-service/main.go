package main

import (
	"context"
	"fmt"
	"joke-api-service/internal/api"
	"joke-api-service/internal/config"
	"log"
	"time"
)

func main() {

	var timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	appConfig := config.Config{JokeURL: "https://official-joke-api.appspot.com/random_joke"}

	client := api.NewClient(appConfig.JokeURL, timeout)

	result, err := client.GetRandomJoke(ctx)

	if err != nil {
		log.Fatalf("Critical error fetching joke: %v", err)
	}

	fmt.Printf("%s, %s", result.Setup, result.Punchline)
}
