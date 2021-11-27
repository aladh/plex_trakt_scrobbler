package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type payload struct {
	Embeds []embed `json:"embeds"`
}

type embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

const contentType = "application/json"

func Track(webhookURL string, err error) error {
	body, err := json.Marshal(payload{
		Embeds: []embed{
			{
				Title:       "Plex Trakt Scrobbler",
				Description: err.Error(),
			},
		},
	})

	resp, err := http.Post(webhookURL, contentType, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error sending webhook: %w", err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("bad response from webhook endpoint: %s", resp.Status)
	}

	return nil
}
