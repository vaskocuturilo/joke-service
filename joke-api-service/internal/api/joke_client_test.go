package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetRandomJoke(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock a successful response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Setup": "Why did the gopher cross the road?", "Punchline":"To get to the other side of the stack!"}`))
	}))

	defer server.Close()

	client := NewClient(server.URL, time.Second*5)

	ctx := context.Background()

	joke, err := client.GetRandomJoke(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if joke.Setup == "" {
		t.Error("Expected a joke value, got empty string")
	}
}

func TestGetRandomJoke_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Id": 1, "Setup": "Why do Gophers hate water?", "Punchline":"Because they prefer the cloud!"}`))
	}))

	defer server.Close()

	client := NewClient(server.URL, time.Second*5)

	ctx := context.Background()

	joke, err := client.GetRandomJoke(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedJokeSetup := "Why do Gophers hate water?"

	expectedJokePunchline := "Because they prefer the cloud!"

	if joke.Setup != expectedJokeSetup {
		t.Errorf("Expected joke %q, got %q", expectedJokeSetup, joke.Setup)
	}

	if joke.Punchline != expectedJokePunchline {
		t.Errorf("Expected joke %q, got %q", expectedJokePunchline, joke.Punchline)
	}
}

func TestGetRandomJoke_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Id": "1", "text": "This JSON is broken"}`))
	}))

	defer server.Close()

	client := NewClient(server.URL, time.Second*5)

	ctx := context.Background()

	_, err := client.GetRandomJoke(ctx)

	if err == nil {
		t.Fatal("expected an error but got nil")
	}

	expectedSubstring := "failed to decode joke"

	if !strings.Contains(err.Error(), expectedSubstring) {
		t.Errorf("Expected error message to contain %q, but got: %v", expectedSubstring, err)
	}
}

func TestGetRandomJoke_ReproduceStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		mockStatus     int
		wantErrContain string
	}{
		{
			name:           "Wrong Status Code",
			mockResponse:   `{}`,
			mockStatus:     http.StatusBadRequest,
			wantErrContain: "unexpected status code: 400",
		},
		{
			name:           "Wrong Status Code",
			mockResponse:   `{}`,
			mockStatus:     http.StatusNotFound,
			wantErrContain: "unexpected status code: 404",
		},
		{
			name:           "Wrong Status Code",
			mockResponse:   `{}`,
			mockStatus:     http.StatusInternalServerError,
			wantErrContain: "unexpected status code: 500",
		},
		{
			name:           "Wrong Status Code",
			mockResponse:   `{}`,
			mockStatus:     http.StatusBadGateway,
			wantErrContain: "unexpected status code: 502",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := NewClient(server.URL, time.Second)
			_, err := client.GetRandomJoke(context.Background())

			if err == nil {
				t.Fatal("expected an error but got nil")
			}

			if !strings.Contains(err.Error(), tt.wantErrContain) {
				t.Errorf("error string %q does not contain %q", err.Error(), tt.wantErrContain)
			}

			t.Logf("Reproduced Error: %v", err)
		})
	}
}
