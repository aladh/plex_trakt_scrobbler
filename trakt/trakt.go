package trakt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type Trakt struct {
	client   *http.Client
	clientID string
}

func New(clientID string, clientSecret string, accessToken string) *Trakt {
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://api.trakt.tv/oauth/token",
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{})
	client := cfg.Client(ctx, &oauth2.Token{TokenType: "Bearer", AccessToken: accessToken})

	return &Trakt{client: client, clientID: clientID}
}

func (t *Trakt) WatchEpisode(ids map[string]string) error {
	reqBody, err := json.Marshal(watchEpisodeRequest(ids, time.Now().UTC()))
	if err != nil {
		return fmt.Errorf("error marshalling trakt request: %w", err)
	}

	log.Printf("Sending watched episode: %s\n", string(reqBody))

	resBody, err := t.watchMedia(reqBody)
	if err != nil {
		return err
	}

	log.Printf("Got response: %s\n", string(resBody))

	return nil
}

func (t *Trakt) WatchMovie(ids map[string]string) error {
	reqBody, err := json.Marshal(watchMovieRequest(ids, time.Now().UTC()))
	if err != nil {
		return fmt.Errorf("error marshalling trakt request: %w", err)
	}

	log.Printf("Sending watched movie: %s\n", string(reqBody))

	resBody, err := t.watchMedia(reqBody)
	if err != nil {
		return err
	}

	log.Printf("Got response: %s\n", string(resBody))

	return nil
}

func (t *Trakt) LatestWatchedMovie() (WatchedMovie, error) {
	log.Println("Getting latest watched movie")

	res, err := t.request("GET", "https://api.trakt.tv/sync/history/movies?limit=1", nil)
	if err != nil {
		return WatchedMovie{}, fmt.Errorf("error making LatestWatchedMovie request: %w", err)
	}

	if res.StatusCode != 200 {
		return WatchedMovie{}, fmt.Errorf("received bad response code %d", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return WatchedMovie{}, fmt.Errorf("error reading trakt response body: %w", err)
	}

	log.Printf("Got response: %s\n", string(resBody))

	var watchedMovies []WatchedMovie

	err = json.Unmarshal(resBody, &watchedMovies)
	if err != nil {
		return WatchedMovie{}, fmt.Errorf("error parsing JSON response: %w", err)
	}

	return watchedMovies[0], nil
}

func (t *Trakt) RemoveFromHistory(id int) error {
	reqBody, err := json.Marshal(removeHistoryRequest(id))
	if err != nil {
		return fmt.Errorf("error marshalling trakt request: %w", err)
	}

	log.Printf("Removing items from history: %s\n", string(reqBody))

	res, err := t.request("POST", "https://api.trakt.tv/sync/history/remove", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("error making RemoveFromHistory request: %w", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("received bad response code %d", res.StatusCode)
	}

	return nil
}

func (t *Trakt) watchMedia(reqBody []byte) ([]byte, error) {
	res, err := t.request("POST", "https://api.trakt.tv/sync/history", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error making watchMedia request: %w", err)
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf("received bad response code %d", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading trakt response body: %w", err)
	}

	return resBody, nil
}

func (t *Trakt) request(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating trakt request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("trakt-api-version", "2")
	req.Header.Set("trakt-api-key", t.clientID)

	res, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to trakt: %w", err)
	}

	return res, nil
}
