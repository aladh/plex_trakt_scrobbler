package trakt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/ali-l/plex_trakt_scrobbler/plex"
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

func (t *Trakt) WatchShow(id plex.ID, season int, episode int) error {
	reqStruct := newWatchShowRequest(id.Value, season, episode, time.Now().UTC())

	reqBody, err := json.Marshal(reqStruct)
	if err != nil {
		return fmt.Errorf("error marshalling trakt request: %w", err)
	}

	log.Printf("Sending watched show: %s\n", string(reqBody))

	req, err := http.NewRequest("POST", "https://api.trakt.tv/sync/history", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("error creating trakt request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("trakt-api-version", "2")
	req.Header.Set("trakt-api-key", t.clientID)

	res, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to trakt: %w", err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading trakt response body: %w", err)
	}

	log.Printf("Got response: %s\n", string(b))

	return nil
}
