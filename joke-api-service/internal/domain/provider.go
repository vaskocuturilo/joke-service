package domain

import (
	"context"
)

type JokeProvider interface {
	GetRandomJoke(ctx context.Context) (Joke, error)
	GetTenJokes(ctx context.Context) ([]Joke, error)
}
