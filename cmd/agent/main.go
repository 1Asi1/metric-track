package main

import (
	"flag"
	"fmt"

	"github.com/1Asi1/metric-track.git/internal/agent/integration"
	"github.com/1Asi1/metric-track.git/internal/agent/service"
	"github.com/1Asi1/metric-track.git/internal/config"
)

var hostServer string
var reportInterval int
var pollInterval int

func parseFlag() {
	flag.StringVar(&hostServer, "a", "localhost:8080", "address and port to run agent")
	flag.IntVar(&reportInterval, "r", 10, "report agent interval")
	flag.IntVar(&pollInterval, "p", 2, "pull agent interval")
	flag.Parse()
}

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
	}

	parseFlag()
	cfg.ReportInterval = reportInterval
	cfg.PollInterval = pollInterval
	cfg.MetricServerAddr = hostServer

	s := service.New(cfg)
	c := integration.New(cfg, s)

	c.SendMetricPeriodic()
}
