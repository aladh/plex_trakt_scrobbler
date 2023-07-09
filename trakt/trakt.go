package trakt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	clientID    string
	accessToken string
	client      http.Client
}

type response struct {
	Body       []byte
	StatusCode int
}

func NewClient(clientID string, accessToken string) *Client {
	return &Client{clientID: clientID, accessToken: accessToken, client: http.Client{}}
}

func (t *Client) WatchEpisode(ids map[string]string) error {
	reqBody, err := json.Marshal(watchEpisodeRequest(ids, time.Now().UTC()))
	if err != nil {
		return fmt.Errorf("error marshalling trakt request: %w", err)
	}

	log.Printf("Sending watched episode: %s\n", string(reqBody))

	res, err := t.request("POST", "https://api.trakt.tv/sync/history", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("error making WatchEpisode request: %w", err)
	}

	if res.StatusCode != 201 {
		return fmt.Errorf("received bad response code %d", res.StatusCode)
	}

	log.Printf("Got response: %s\n", string(res.Body))

	if strings.Contains(string(res.Body), "ids") {
		return fmt.Errorf("show not found")
	}

	return nil
}

func (t *Client) WatchMovie(ids map[string]string) error {
	reqBody, err := json.Marshal(watchMovieRequest(ids, time.Now().UTC()))
	if err != nil {
		return fmt.Errorf("error marshalling trakt request: %w", err)
	}

	log.Printf("Sending watched movie: %s\n", string(reqBody))

	res, err := t.request("POST", "https://api.trakt.tv/sync/history", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("error making WatchMovie request: %w", err)
	}

	if res.StatusCode != 201 {
		return fmt.Errorf("received bad response code %d", res.StatusCode)
	}

	log.Printf("Got response: %s\n", string(res.Body))

	if strings.Contains(string(res.Body), "ids") {
		return fmt.Errorf("movie not found")
	}

	return nil
}

func (t *Client) LatestWatchedMovies() ([]WatchedMovie, error) {
	log.Println("Getting latest watched movie")

	res, err := t.request("GET", "https://api.trakt.tv/sync/history/movies?limit=5", nil)
	if err != nil {
		return []WatchedMovie{}, fmt.Errorf("error making LatestWatchedMovies request: %w", err)
	}

	if res.StatusCode != 200 {
		return []WatchedMovie{}, fmt.Errorf("received bad response code %d", res.StatusCode)
	}

	log.Printf("Got response: %s\n", string(res.Body))

	var watchedMovies []WatchedMovie

	err = json.Unmarshal(res.Body, &watchedMovies)
	if err != nil {
		return []WatchedMovie{}, fmt.Errorf("error parsing JSON response: %w", err)
	}

	return watchedMovies, nil
}

func (t *Client) RemoveFromHistory(id int) error {
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

func (t *Client) request(method string, url string, body io.Reader) (response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return response{}, fmt.Errorf("error creating trakt request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("trakt-api-version", "2")
	req.Header.Set("trakt-api-key", t.clientID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.accessToken))

	res, err := t.client.Do(req)
	if err != nil {
		return response{}, fmt.Errorf("error sending request to trakt: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return response{}, fmt.Errorf("error reading trakt response body: %w", err)
	}

	return response{Body: resBody, StatusCode: res.StatusCode}, nil
}
