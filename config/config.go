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
}

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

	return &Config{
		TraktClientID:     traktClientID,
		TraktClientSecret: traktClientSecret,
		TraktAccessToken:  traktAccessToken,
		TraktRefreshToken: traktRefreshToken,
	}, nil
}

func getEnvString(name string) (string, error) {
	value, exists := os.LookupEnv(name)
	if !exists {
		return "", fmt.Errorf("environment variable %s not found", name)
	}

	return value, nil
}
