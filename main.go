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

		if !plex.ShouldProcess(payload, cfg.PlexServerUUID) {
			return
		}

		log.Printf("Received scrobble: %v\n", payload)

		err = traktClient.WatchEpisode(payload.Metadata.ID().Value, payload.Metadata.Season(), payload.Metadata.Episode())
		if err != nil {
			log.Printf("Error watching episode: %s\n", err)
			return
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
