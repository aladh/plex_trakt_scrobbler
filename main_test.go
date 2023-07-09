package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aladh/plex_trakt_scrobbler/config"
	"github.com/aladh/plex_trakt_scrobbler/trakt"
)

func TestWatchMovieWebhook(t *testing.T) {
	_, enabled := os.LookupEnv("ENABLE_E2E_TESTS")
	if !enabled {
		t.Skip("End to end tests are not enabled")
	}

	// Start server and wait for it to be ready
	go main()
	time.Sleep(1 * time.Second)

	cfg, err := config.FromEnv()
	if err != nil {
		t.Fatalf("error loading config from env: %s", err)
	}

	sendWebhook(t, cfg)
	verifyAndCleanup(t, cfg)
}

func verifyAndCleanup(t *testing.T, cfg *config.Config) {
	traktClient := trakt.NewClient(cfg.TraktClientID, cfg.TraktAccessToken)
	watchedMovies, err := traktClient.LatestWatchedMovies()
	if err != nil {
		t.Fatalf("error getting latest watchedMovies: %s", err)
	}

	expectedTitle := "McHale's Navy"

	var expectedMovie trakt.WatchedMovie
	for _, movie := range watchedMovies {
		if movie.Movie.Title == expectedTitle {
			expectedMovie = movie
			break
		}
	}

	if expectedMovie.Movie.Title != expectedTitle {
		t.Fatalf("title = %s, want %s", expectedMovie.Movie.Title, expectedTitle)
	}

	// Comply with rate limiting
	time.Sleep(1 * time.Second)

	err = traktClient.RemoveFromHistory(expectedMovie.ID)
	if err != nil {
		t.Fatalf("error removing from history: %s", err)
	}
}

func sendWebhook(t *testing.T, cfg *config.Config) {
	payload, err := os.ReadFile("testdata/webhook.json")
	if err != nil {
		t.Fatalf("error opening fixture file: %s", err)
	}

	var formBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&formBuffer)
	payloadWriter, err := multipartWriter.CreateFormField("payload")
	if err != nil {
		t.Fatalf("error creating form field: %s", err)
	}

	_, err = payloadWriter.Write(payload)
	if err != nil {
		t.Fatalf("error writing form field: %s", err)
	}

	err = multipartWriter.Close()
	if err != nil {
		t.Fatalf("error closing form writer: %s", err)
	}

	response, err := http.Post(fmt.Sprintf("http://localhost:%s", cfg.Port), multipartWriter.FormDataContentType(), &formBuffer)
	if err != nil {
		t.Fatalf("error making request to app server: %s", err)
	}

	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("bad response from app server")
	}
}
