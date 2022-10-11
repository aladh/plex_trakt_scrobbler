package notifier

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func NotifyMovieScrobble(webhookURL string) error {
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
