package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	PollInterval     time.Duration
	ReportInterval   time.Duration
	MetricServerAddr string
}

func New() (Config, error) {
	var cfg Config

	add := flag.String("a", "localhost:8080", "address and port to run agent")
	rep := flag.Int("r", 10, "report agent interval")
	pull := flag.Int("p", 2, "pull agent interval")
	flag.Parse()

	pollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		pI, err := strconv.Atoi(pollInterval)
		if err != nil {
			return Config{}, err
		}

		cfg.PollInterval = time.Duration(pI) * time.Second
	} else {
		cfg.PollInterval = time.Duration(*pull) * time.Second
	}

	reportInterval, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		rI, err := strconv.Atoi(reportInterval)
		if err != nil {
			return Config{}, err
		}

		cfg.ReportInterval = time.Duration(rI) * time.Second
	} else {
		cfg.ReportInterval = time.Duration(*rep) * time.Second
	}

	metricServerAddr, ok := os.LookupEnv("ADDRESS")
	if ok {
		cfg.MetricServerAddr = metricServerAddr
	} else {
		cfg.MetricServerAddr = *add
	}

	return cfg, nil
}
