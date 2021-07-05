package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ali-l/plex_trakt_scrobbler/config"
	"github.com/ali-l/plex_trakt_scrobbler/server"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("Error getting config from env: %s", err)
	}

	http.HandleFunc("/", server.Handler(cfg))

	log.Printf("Starting server on port %s\n", cfg.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil)
	if err != nil {
		log.Fatalf("Error starting web server: %s", err)
	}
}
