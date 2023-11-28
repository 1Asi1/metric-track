package main

import (
	"os"

	"github.com/1Asi1/metric-track.git/internal/agent/config"
	"github.com/1Asi1/metric-track.git/internal/agent/integration"
	"github.com/1Asi1/metric-track.git/internal/agent/service"
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

	s := service.New(cfg, l)
	c := integration.New(cfg, s, l)

	c.SendMetricPeriodic()
}
