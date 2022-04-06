package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	TraktClientID    string
	TraktAccessToken string
	PlexServerUUIDs  []string
	PlexUsername     string
	Port             string
	WebhookURL       string
}

const DefaultPort = "8080"

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
		port = DefaultPort
	}

	webhookURL, err := getEnvString("WEBHOOK_URL")
	if err != nil {
		return nil, err
	}

	return &Config{
		TraktClientID:    traktClientID,
		TraktAccessToken: traktAccessToken,
		PlexServerUUIDs:  plexServerUUIDs,
		PlexUsername:     plexUsername,
		Port:             port,
		WebhookURL:       webhookURL,
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
