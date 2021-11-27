package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aladh/plex_trakt_scrobbler/config"
	"github.com/aladh/plex_trakt_scrobbler/errors"
	"github.com/aladh/plex_trakt_scrobbler/plex"
	"github.com/aladh/plex_trakt_scrobbler/trakt"
)

const postMethod = "POST"

func Handler(cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	traktClient := trakt.New(cfg.TraktClientID, cfg.TraktClientSecret, cfg.TraktAccessToken)

	return func(_ http.ResponseWriter, request *http.Request) {
		err := processRequest(cfg, traktClient, request)
		if err != nil {
			log.Println(err)
			log.Printf("Request payload: %s\n", request.FormValue("payload"))

			err = errors.Track(cfg.WebhookURL, err)
			if err != nil {
				log.Printf("error tracking error: %s\n", err)
			}
		}
	}
}

func processRequest(cfg *config.Config, traktClient *trakt.Trakt, request *http.Request) error {
	// Plex webhooks are always POST
	if request.Method != postMethod {
		return nil
	}

	payload, err := parsePayload(request)
	if err != nil {
		return fmt.Errorf("error parsing webhook payload: %w", err)
	}

	if !plex.ShouldProcess(payload, cfg.PlexServerUUIDs, cfg.PlexUsername) {
		return nil
	}

	log.Printf("Received scrobble: %v\n", payload)

	err = watchMedia(payload, traktClient)
	if err != nil {
		return fmt.Errorf("error watching media: %w", err)
	}

	return nil
}

func parsePayload(request *http.Request) (*plex.Payload, error) {
	err := request.ParseMultipartForm(10_000)
	if err != nil {
		return nil, fmt.Errorf("error parsing multipart form: %w", err)
	}

	var payload plex.Payload

	err = json.Unmarshal([]byte(request.FormValue("payload")), &payload)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling webhook payload: %w", err)
	}

	return &payload, nil
}

func watchMedia(payload *plex.Payload, traktClient *trakt.Trakt) error {
	var err error

	switch payload.Type() {
	case plex.ShowType:
		err = traktClient.WatchEpisode(payload.IDs())
	case plex.MovieType:
		err = traktClient.WatchMovie(payload.IDs())
	default:
		err = fmt.Errorf("unrecognized media type %s", payload.Type())
	}

	if err != nil {
		return err
	}

	return nil
}
