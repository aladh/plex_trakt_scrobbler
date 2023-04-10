package notifier

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aladh/plex_trakt_scrobbler/config"
	"github.com/aladh/plex_trakt_scrobbler/contextkeys"
	"github.com/aladh/plex_trakt_scrobbler/plex"
)

func NotifyScrobble(ctx context.Context) error {
	cfg := ctx.Value(contextkeys.Config).(*config.Config)
	payload := ctx.Value(contextkeys.Payload).(*plex.Payload)

	if payload.Type() == plex.MovieType {
		return notifyMovieScrobble(cfg.MovieScrobbleWebhookURL)
	}

	return nil
}

func notifyMovieScrobble(webhookURL string) error {
	if len(webhookURL) == 0 {
		log.Println("no movie scrobble URL specified")
		return nil
	}

	resp, err := http.Post(webhookURL, "text/plain", strings.NewReader(""))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("bad response from endpoint: %s", resp.Status)
	}

	return nil
}
