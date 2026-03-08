package main

import (
	"log"

	cfg "org-structure-api/internal/config"
	db "org-structure-api/internal/database"
)

func main() {
	cfg, err := cfg.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pg, err := db.NewPostgresClient(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	srv, err := container(cfg.Server, pg)
	if err != nil {
		log.Fatalf("failed to build container: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
