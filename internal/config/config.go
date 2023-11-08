package config

import (
	"os"
	"strconv"
)

type Config struct {
	PollInterval     int
	ReportInterval   int
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

	metricServerAddr := os.Getenv("")
	if metricServerAddr == "" {
		metricServerAddr = "http://localhost:8080"
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
			MetricServerAddr: metricServerAddr,
		},
		nil
}
