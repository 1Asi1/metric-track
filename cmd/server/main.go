package main

import (
	"os"

	"github.com/1Asi1/metric-track.git/internal/config"
	"github.com/1Asi1/metric-track.git/internal/server/apiserver"
	"github.com/rs/zerolog"
)

func main() {
	out := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05 -0700",
		NoColor:    true,
	}

	l := zerolog.New(out)

	l = l.Level(zerolog.InfoLevel).With().Timestamp().Logger()

	cfg, err := config.New(l)
	if err != nil {
		l.Fatal().Err(err).Msg("config.New")
	}

	server := apiserver.New(cfg, l)
	if err = server.Run(); err != nil {
		l.Fatal().Err(err).Msg("http.ListenAndServe")
	}
}
