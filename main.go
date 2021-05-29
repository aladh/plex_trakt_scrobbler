package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ali-l/plex_trakt_scrobbler/config"
	"github.com/ali-l/plex_trakt_scrobbler/plex"
	"github.com/ali-l/plex_trakt_scrobbler/trakt"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("Error getting config from env: %s", err)
	}

	traktClient := trakt.New(cfg.TraktClientID, cfg.TraktClientSecret, cfg.TraktAccessToken, cfg.TraktRefreshToken)

	http.HandleFunc("/", handler(cfg, traktClient))

	log.Printf("Starting server on port %s\n", cfg.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil)
	if err != nil {
		log.Fatalf("Error starting web server: %s", err)
	}
}

func handler(cfg *config.Config, traktClient *trakt.Trakt) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Plex webhooks are always POST
		if request.Method != "POST" {
			return
		}

		payload, err := parsePayload(request)
		if err != nil {
			log.Printf("Error parsing webhook payload: %s\n", err)
			return
		}

		if !plex.ShouldProcess(payload, cfg.PlexServerUUIDs, cfg.PlexUsername) {
			return
		}

		log.Printf("Received scrobble: %v\n", payload)

		err = watchMedia(payload, traktClient)
		if err != nil {
			log.Printf("Error watching media: %s\n", err)
		}
	}
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
