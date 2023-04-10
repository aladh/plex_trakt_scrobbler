package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aladh/plex_trakt_scrobbler/config"
	contextkeys "github.com/aladh/plex_trakt_scrobbler/context"
	"github.com/aladh/plex_trakt_scrobbler/errortracker"
	"github.com/aladh/plex_trakt_scrobbler/notifier"
	"github.com/aladh/plex_trakt_scrobbler/plex"
	"github.com/aladh/plex_trakt_scrobbler/trakt"
	"github.com/aladh/plex_trakt_scrobbler/util"
)

const postMethod = "POST"

func Handler(cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	traktClient := trakt.NewClient(cfg.TraktClientID, cfg.TraktAccessToken)

	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.WithValue(request.Context(), contextkeys.Config, cfg)
		ctx = context.WithValue(ctx, contextkeys.TraktClient, traktClient)
		ctx = context.WithValue(ctx, contextkeys.Request, request)

		// Plex webhooks are always POST so we can ignore other methods
		if request.Method != postMethod {
			log.Printf("Invalid request method: %s, expected POST", request.Method)
			http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		err := processRequest(ctx)
		if err != nil {
			log.Println(err)
			log.Printf("Request payload: %s\n", request.FormValue("payload"))

			err = errortracker.Track(cfg.ErrorWebhookURL, err)
			if err != nil {
				log.Printf("error tracking error: %s\n", err)
			}

			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}

func processRequest(ctx context.Context) error {
	request := ctx.Value(contextkeys.Request).(*http.Request)

	payload, err := parsePayload(request)
	if err != nil {
		return fmt.Errorf("error parsing webhook payload: %w", err)
	}
	ctx = context.WithValue(ctx, contextkeys.Payload, payload)

	if !isAuthorized(ctx) {
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

	err = processScrobble(ctx)
	if err != nil {
		return fmt.Errorf("error processing scrobble: %w", err)
	}

	err = notifier.NotifyScrobble(ctx)
	if err != nil {
		return fmt.Errorf("error notifying scrobble: %w", err)
	}

	return nil
}

func isAuthorized(ctx context.Context) bool {
	payload := ctx.Value(contextkeys.Payload).(*plex.Payload)
	cfg := ctx.Value(contextkeys.Config).(*config.Config)

	// Check that the webhook is coming from an allowed server
	if !util.SliceContains(cfg.PlexServerUUIDs, payload.ServerUUID()) {
		log.Printf("Unauthorized request from server UUID: %s\n", payload.ServerUUID())
		return false
	}

	// Only scrobble plays from the specified user
	if payload.Username() != cfg.PlexUsername {
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

func processScrobble(ctx context.Context) error {
	payload := ctx.Value(contextkeys.Payload).(*plex.Payload)
	traktClient := ctx.Value(contextkeys.TraktClient).(*trakt.Client)

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
