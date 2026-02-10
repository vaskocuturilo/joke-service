package api

import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"joke-api-service/internal/domain"
	"log"
	"net/http"
	"net/url"
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
	var result domain.Joke

	err := j.decodeResponse(ctx, "/random_joke", &result)

	return result, err
}

func (j *Client) GetTenJokes(ctx context.Context) ([]domain.Joke, error) {
	var result []domain.Joke

	err := j.decodeResponse(ctx, "/random_ten", &result)

	return result, err
}

func (j *Client) decodeResponse(ctx context.Context, endpoint string, target any) error {
	fullURL, err := url.JoinPath(j.baseURL, endpoint)

	if err != nil {
		return fmt.Errorf("failed to create full path: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	response, err := j.httpClient.Do(request)

	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close body: %v", err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode joke: %w", err)
	}
	return nil
}
