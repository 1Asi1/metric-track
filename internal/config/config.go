package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	PollInterval     time.Duration
	ReportInterval   time.Duration
	MetricServerAddr string
}

func New(log zerolog.Logger) (Config, error) {
	l := log.With().Str("config", "New").Logger()

	var cfg Config

	add := flag.String("a", "localhost:8080", "address and port to run agent")
	rep := flag.Int("r", 10, "report agent interval")
	pull := flag.Int("p", 2, "pull agent interval")
	flag.Parse()

	pollInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		pI, err := strconv.Atoi(pollInterval)
		if err != nil {
			l.Error().Err(err).Msgf("strconv.Atoi, poll interval value: %s", pollInterval)
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
			l.Error().Err(err).Msgf("strconv.Atoi, report interval value: %s", reportInterval)
			return Config{}, err
		}

		cfg.ReportInterval = time.Duration(rI) * time.Second
	} else {
		cfg.ReportInterval = time.Duration(*rep) * time.Second
	}

	metricServerAddr, ok := os.LookupEnv("ADDRESS")
	if ok {
		l.Info().Msgf("server address value: %s", metricServerAddr)
		cfg.MetricServerAddr = metricServerAddr
	} else {
		l.Info().Msgf("server address value: %s", *add)
		cfg.MetricServerAddr = *add
	}

	return cfg, nil
}
