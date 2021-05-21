package plex

import (
	"log"
	"strings"
)

const scrobbleEvent = "media.scrobble"

func ShouldProcess(payload *Payload, allowedUUIDs string) bool {
	// Check that the webhook is coming from the right server
	if !strings.Contains(allowedUUIDs, payload.Server.UUID) {
		log.Printf("Unauthorized request from server UUID: %s\n", payload.Server.UUID)
		return false
	}

	// Only send watch request when media has been completely watched
	if payload.Event != scrobbleEvent {
		return false
	}

	return true
}
