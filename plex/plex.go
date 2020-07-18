package plex

import "log"

const scrobbleEvent = "media.scrobble"

func ShouldProcess(payload *Payload, serverUUID string) bool {
	// Check that the webhook is coming from the right server
	if payload.Server.UUID != serverUUID {
		log.Printf("Unauthorized request from server: %s\n", payload.Server.UUID)
		return false
	}

	// Only send watch request when media has been completely watched
	if payload.Event != scrobbleEvent {
		return false
	}

	return true
}
