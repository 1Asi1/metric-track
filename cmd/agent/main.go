package main

import (
	"fmt"
	"log"

	"github.com/1Asi1/metric-track.git/internal/agent/config"
	"github.com/1Asi1/metric-track.git/internal/agent/integration"
	"github.com/1Asi1/metric-track.git/internal/agent/service"
	"github.com/1Asi1/metric-track.git/internal/logger"
	"github.com/rs/zerolog"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)

	cfg, err := config.New(logger.NewLogger())
	if err != nil {
		log.Fatal("config.New")
	}

	l := logger.NewLogger()
	l = l.Level(zerolog.InfoLevel).With().Timestamp().Logger()

	s := service.New(cfg, l)
	c := integration.New(cfg, s, l)

	c.SendMetricPeriodic()
}
