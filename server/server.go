package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aladh/plex_trakt_scrobbler/config"
	"github.com/aladh/plex_trakt_scrobbler/errors"
	"github.com/aladh/plex_trakt_scrobbler/notifier"
	"github.com/aladh/plex_trakt_scrobbler/plex"
	"github.com/aladh/plex_trakt_scrobbler/trakt"
	"github.com/aladh/plex_trakt_scrobbler/util"
)

const postMethod = "POST"

func Handler(cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	traktClient := trakt.New(cfg.TraktClientID, cfg.TraktAccessToken)

	return func(w http.ResponseWriter, request *http.Request) {
		err := processRequest(cfg, traktClient, request)
		if err != nil {
			log.Println(err)
			log.Printf("Request payload: %s\n", request.FormValue("payload"))

			err = errors.Track(cfg.ErrorWebhookURL, err)
			if err != nil {
				log.Printf("error tracking error: %s\n", err)
			}

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
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

	if !isAuthorized(payload, cfg.PlexServerUUIDs, cfg.PlexUsername) {
		return nil
	}

	// Only send watch request when media has been completely watched
	if !payload.IsScrobble() {
		log.Printf("Skipping non-scrobble event: %s\n", payload.Event)
		return nil
	}

	log.Printf("Parsed payload from scrobble: %v\n", payload)

	if !payload.HasIDs() {
		return fmt.Errorf("error processing request for title %s: payload has no IDs", payload.Metadata.Title)
	}

	err = watchMedia(payload, traktClient, cfg)
	if err != nil {
		return fmt.Errorf("error watching media: %w", err)
	}

	return nil
}

func isAuthorized(payload *plex.Payload, allowedUUIDs []string, allowedUsername string) bool {
	// Check that the webhook is coming from an allowed server
	if !util.SliceContains(allowedUUIDs, payload.ServerUUID()) {
		log.Printf("Unauthorized request from server UUID: %s\n", payload.ServerUUID())
		return false
	}

	// Only scrobble plays from the specified user
	if payload.Username() != allowedUsername {
		log.Printf("User not recognized: %s\n", payload.Username())
		return false
	}

	return true
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

func watchMedia(payload *plex.Payload, traktClient *trakt.Trakt, cfg *config.Config) error {
	var err error

	switch payload.Type() {
	case plex.ShowType:
		err = traktClient.WatchEpisode(payload.IDs())
	case plex.MovieType:
		err = traktClient.WatchMovie(payload.IDs())
		if err != nil {
			return err
		}

		err = notifier.NotifyMovieScrobble(cfg.MovieScrobbleWebhookURL)
	default:
		err = fmt.Errorf("unrecognized media type %s", payload.Type())
	}

	if err != nil {
		return err
	}

	return nil
}
