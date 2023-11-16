package main

import (
	"log"

	"github.com/1Asi1/metric-track.git/internal/agent/integration"
	"github.com/1Asi1/metric-track.git/internal/agent/service"
	"github.com/1Asi1/metric-track.git/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	s := service.New(cfg)
	c := integration.New(cfg, s)

	c.SendMetricPeriodic()
}
