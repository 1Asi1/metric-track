package main

import (
	"log"

	"github.com/1Asi1/metric-track.git/internal/config"
	"github.com/1Asi1/metric-track.git/internal/server/apiserver"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	server := apiserver.New(cfg)
	if err = server.Run(); err != nil {
		log.Fatalf("http.ListenAndServe panic: %v", err)
	}
}
