package main

import (
	"log"

	"github.com/1Asi1/metric-track.git/internal/logger"
	"github.com/1Asi1/metric-track.git/internal/server/apiserver"
	"github.com/1Asi1/metric-track.git/internal/server/config"
	"github.com/rs/zerolog"
)

func main() {
	cfg, err := config.New(logger.NewLogger())
	if err != nil {
		log.Fatal("config.New")
	}

	l := logger.NewLogger()
	l = l.Level(zerolog.InfoLevel).With().Timestamp().Logger()

	server := apiserver.New(cfg, l)
	if err = server.Run(); err != nil {
		l.Fatal().Err(err).Msg("server.Run")
	}
}
