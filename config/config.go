package config

import (
	"fmt"
	"os"
)

type Config struct {
	TraktClientID     string
	TraktClientSecret string
	TraktAccessToken  string
	TraktRefreshToken string
	PlexServerUUIDs   string
	PlexUsername      string
	Port              string
}

const defaultPort = "8080"

func FromEnv() (*Config, error) {
	traktClientID, err := getEnvString("TRAKT_CLIENT_ID")
	if err != nil {
		return nil, err
	}

	traktClientSecret, err := getEnvString("TRAKT_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}

	traktAccessToken, err := getEnvString("TRAKT_ACCESS_TOKEN")
	if err != nil {
		return nil, err
	}

	traktRefreshToken, err := getEnvString("TRAKT_REFRESH_TOKEN")
	if err != nil {
		return nil, err
	}

	plexServerUUIDs, err := getEnvString("PLEX_SERVER_UUIDS")
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

	return &Config{
		TraktClientID:     traktClientID,
		TraktClientSecret: traktClientSecret,
		TraktAccessToken:  traktAccessToken,
		TraktRefreshToken: traktRefreshToken,
		PlexServerUUIDs:   plexServerUUIDs,
		PlexUsername:      plexUsername,
		Port:              port,
	}, nil
}

func getEnvString(name string) (string, error) {
	value, exists := os.LookupEnv(name)
	if !exists {
		return "", fmt.Errorf("environment variable %s not found", name)
	}

	return value, nil
}
