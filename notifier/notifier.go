package notifier

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aladh/plex_trakt_scrobbler/config"
	"github.com/aladh/plex_trakt_scrobbler/plex"
)

func NotifyScrobble(cfg *config.Config, payloadType plex.PayloadType) error {
	if payloadType == plex.MovieType {
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

	if resp.StatusCode != 201 {
		return fmt.Errorf("bad response from endpoint: %s", resp.Status)
	}

	return nil
}
