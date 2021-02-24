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

func New(clientID string, clientSecret string, accessToken string, refreshToken string) *Trakt {
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://api.trakt.tv/oauth/token",
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{})
	client := cfg.Client(ctx, &oauth2.Token{TokenType: "Bearer", AccessToken: accessToken, RefreshToken: refreshToken})

	return &Trakt{client: client, clientID: clientID}
}

func (t *Trakt) WatchEpisode(id string, season int, episode int) error {
	reqBody, err := json.Marshal(watchEpisodeRequest(id, season, episode, time.Now().UTC()))
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

func (t *Trakt) WatchMovie(id string) error {
	reqBody, err := json.Marshal(watchMovieRequest(id, time.Now().UTC()))
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

func (t *Trakt) watchMedia(reqBody []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", "https://api.trakt.tv/sync/history", bytes.NewReader(reqBody))
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

	if res.StatusCode != 201 {
		return nil, fmt.Errorf("received bad response code %d", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading trakt response body: %w", err)
	}

	return resBody, nil
}
