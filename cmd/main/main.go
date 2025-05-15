package main

import (
	"goproject/internal/config"
	"goproject/internal/storage/postgres"
	"log"
	"net/http"
)

func main() {
	cfg, msg := config.MustLoad()

	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer storage.Close()

	log.Println(msg)

	http.ListenAndServe(cfg.HTTPServer.Address, nil)

}
