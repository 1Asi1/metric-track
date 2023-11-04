package config

import (
	"os"
	"strconv"
)

type Config struct {
	PollInterval     int
	ReportInterval   int
	MetricServerPort string
	MetricServerAddr string
}

func New() (Config, error) {
	pollInterval := os.Getenv("")
	if pollInterval == "" {
		pollInterval = "2"
	}

	reportInterval := os.Getenv("")
	if reportInterval == "" {
		reportInterval = "10"
	}

	metricServerPort := os.Getenv("")
	if metricServerPort == "" {
		metricServerPort = "8080"
	}

	metricServerAddr := os.Getenv("")
	if metricServerAddr == "" {
		metricServerAddr = "http://localhost:" + metricServerPort
	}

	pI, err := strconv.Atoi(pollInterval)
	if err != nil {
		return Config{}, err
	}

	rI, err := strconv.Atoi(reportInterval)
	if err != nil {
		return Config{}, err
	}

	return Config{
			PollInterval:     pI,
			ReportInterval:   rI,
			MetricServerPort: ":" + metricServerPort,
			MetricServerAddr: metricServerAddr,
		},
		nil
}
