package plex

import "log"

func ShouldProcess(payload *Payload, serverUUID string) bool {
	// Check that the webhook is coming from the right server
	if payload.Server.UUID != serverUUID {
		log.Printf("Unauthorized request from server: %s\n", payload.Server.UUID)
		return false
	}

	// Only send watch request when media has been completely watched
	if payload.Event != "media.scrobble" {
		return false
	}

	// Can only handle shows right now
	if payload.Metadata.LibrarySectionType != "show" {
		return false
	}

	return true
}
