package config

import (
	"fmt"
	"os"
	"strings"
)

const defaultPort = "8080"

type Config struct {
	TraktClientID           string
	TraktAccessToken        string
	PlexServerUUIDs         []string
	PlexUsername            string
	Port                    string
	ErrorWebhookURL         string
	MovieScrobbleWebhookURL string
}

func FromEnv() (*Config, error) {
	traktClientID, err := getEnvString("TRAKT_CLIENT_ID")
	if err != nil {
		return nil, err
	}

	traktAccessToken, err := getEnvString("TRAKT_ACCESS_TOKEN")
	if err != nil {
		return nil, err
	}

	plexServerUUIDs, err := getEnvSlice("PLEX_SERVER_UUIDS")
	if err != nil {
		return nil, err
	}

	plexUsername, err := getEnvString("PLEX_USERNAME")
	if err != nil {
		return nil, err
	}

	port, err := getEnvString("PORT")
	if err != nil {
		port = defaultPort
	}

	errorWebhookURL, err := getEnvString("ERROR_WEBHOOK_URL")
	if err != nil {
		return nil, err
	}

	movieScrobbleWebhookURL, err := getEnvString("MOVIE_SCROBBLE_WEBHOOK_URL")
	if err != nil {
		movieScrobbleWebhookURL = ""
	}

	return &Config{
		TraktClientID:           traktClientID,
		TraktAccessToken:        traktAccessToken,
		PlexServerUUIDs:         plexServerUUIDs,
		PlexUsername:            plexUsername,
		Port:                    port,
		ErrorWebhookURL:         errorWebhookURL,
		MovieScrobbleWebhookURL: movieScrobbleWebhookURL,
	}, nil
}

func getEnvSlice(name string) ([]string, error) {
	const separator = "|"

	value, exists := os.LookupEnv(name)
	if !exists {
		return nil, fmt.Errorf("environment variable %s not found", name)
	}

	return strings.Split(value, separator), nil
}

func getEnvString(name string) (string, error) {
	value, exists := os.LookupEnv(name)
	if !exists {
		return "", fmt.Errorf("environment variable %s not found", name)
	}

	return value, nil
}
