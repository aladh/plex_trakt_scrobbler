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

	// Start server
	go main()

	sendRequest(t)

	cleanup(t)
}

func cleanup(t *testing.T) {
	cfg, err := config.FromEnv()
	if err != nil {
		t.Fatalf("error loading config from env: %s", err)
	}

	traktClient := trakt.New(cfg.TraktClientID, cfg.TraktClientSecret, cfg.TraktAccessToken)
	watchedMovie, err := traktClient.LatestWatchedMovie()
	if err != nil {
		t.Fatalf("error getting latest watchedMovie: %s", err)
	}

	expectedTitle := "McHale's Navy"
	if watchedMovie.Movie.Title != expectedTitle {
		t.Fatalf("title = %s, want %s", watchedMovie.Movie.Title, expectedTitle)
	}

	// Comply with rate limiting
	time.Sleep(1 * time.Second)

	err = traktClient.RemoveFromHistory(watchedMovie.ID)
	if err != nil {
		t.Fatalf("error removing from history: %s", err)
	}
}

func sendRequest(t *testing.T) {
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

	response, err := http.Post(fmt.Sprintf("http://localhost:%s", config.DefaultPort), multipartWriter.FormDataContentType(), &formBuffer)
	if err != nil {
		t.Fatalf("error making request to app server: %s", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("bad response from app server")
	}
}
