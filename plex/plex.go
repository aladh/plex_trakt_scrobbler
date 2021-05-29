package plex

import (
	"log"
	"strings"
)

const scrobbleEvent = "media.scrobble"

func ShouldProcess(payload *Payload, allowedUUIDs string, allowedUsername string) bool {
	// Check that the webhook is coming from the right server
	if !strings.Contains(allowedUUIDs, payload.Server.UUID) {
		log.Printf("Unauthorized request from server UUID: %s\n", payload.Server.UUID)
		return false
	}

	// Only send watch request when media has been completely watched
	if payload.Event != scrobbleEvent {
		return false
	}

	// Only scrobble plays from the specified user
	if payload.Account.Title != allowedUsername {
		log.Printf("User not recognized: %s\n", payload.Account.Title)
		return false
	}

	return true
}
