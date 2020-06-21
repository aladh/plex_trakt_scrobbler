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

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// Plex webhooks are always POST
		if request.Method != "POST" {
			return
		}

		err := request.ParseMultipartForm(10_000)
		if err != nil {
			log.Printf("Error parsing webhook form: %s\n", err)
			return
		}

		var payload plex.Payload

		err = json.Unmarshal([]byte(request.FormValue("payload")), &payload)
		if err != nil {
			log.Printf("Error unmarshaling webhook body: %s\n", err)
			return
		}

		if !plex.ShouldProcess(payload, cfg.PlexServerUUID) {
			return
		}

		log.Printf("Received scrobble: %v\n", payload)

		id, err := payload.Metadata.ID()
		if err != nil {
			log.Printf("Error parsing ID: %s\n", err)
			return
		}

		err = traktClient.WatchShow(id.Value, payload.Metadata.Season(), payload.Metadata.Episode())
		if err != nil {
			log.Printf("Error watching show: %s\n", err)
			return
		}
	})

	log.Printf("Starting server on port %s\n", cfg.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil)
	if err != nil {
		log.Fatalf("Error starting web server: %s", err)
	}
}
