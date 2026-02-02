package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"joke-api-service/internal/domain"
	"log"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (j *Client) GetRandomJoke(ctx context.Context) (domain.Joke, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", j.baseURL, nil)

	if err != nil {
		return domain.Joke{}, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := j.httpClient.Do(request)

	if err != nil {
		return domain.Joke{}, fmt.Errorf("failed to execute request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close body: %v", err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return domain.Joke{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var result domain.Joke

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return domain.Joke{}, fmt.Errorf("failed to decode joke: %w", err)
	}
	return result, nil
}
