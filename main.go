package main

import (
	"encoding/json"
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
		}

		var payload plex.Payload

		err = json.Unmarshal([]byte(request.FormValue("payload")), &payload)
		if err != nil {
			log.Printf("Error unmarshaling webhook body: %s\n", err)
		}

		// Check that the webhook is coming from the right server
		if payload.Server.UUID != cfg.PlexServerUUID {
			log.Printf("Unauthorized request from server: %s\n", payload.Server.UUID)
			return
		}

		// Only send watch request when media has been completely watched
		if payload.Event != "media.scrobble" {
			return
		}

		// Can only handle shows right now
		if payload.Metadata.LibrarySectionType != "show" {
			return
		}

		log.Printf("Received scrobble: %v\n", payload)

		err = traktClient.WatchShow(payload.Metadata.ID(), payload.Metadata.Season(), payload.Metadata.Episode())
		if err != nil {
			log.Printf("Error watching show: %s\n", err)
		}
	})

	err = http.ListenAndServe(":5678", nil)
	if err != nil {
		log.Fatalf("Error starting web server: %s", err)
	}
}
