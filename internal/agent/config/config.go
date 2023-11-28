package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	MetricServerAddr string
	PollInterval     time.Duration
	ReportInterval   time.Duration
}

func New(log zerolog.Logger) (Config, error) {
	l := log.With().Str("config", "New").Logger()

	var cfg Config

	add := flag.String("a", "localhost:8080", "address and port to run agent")
	rep := flag.Int("r", 10, "report agent interval")
	pull := flag.Int("p", 2, "pull agent interval")
	flag.Parse()

	metricServerAddrEnv, ok := os.LookupEnv("ADDRESS")
	if ok {
		l.Info().Msgf("server address value: %s", metricServerAddrEnv)
		cfg.MetricServerAddr = metricServerAddrEnv
	} else {
		l.Info().Msgf("server address value: %s", *add)
		cfg.MetricServerAddr = *add
	}

	pollIntervalEnv, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		pI, err := strconv.Atoi(pollIntervalEnv)
		if err != nil {
			l.Error().Err(err).Msgf("strconv.Atoi, poll interval value: %s", pollIntervalEnv)
			return Config{}, err
		}

		cfg.PollInterval = time.Duration(pI) * time.Second
	} else {
		cfg.PollInterval = time.Duration(*pull) * time.Second
	}

	reportIntervalEnv, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		rI, err := strconv.Atoi(reportIntervalEnv)
		if err != nil {
			l.Error().Err(err).Msgf("strconv.Atoi, report interval value: %s", reportIntervalEnv)
			return Config{}, err
		}

		cfg.ReportInterval = time.Duration(rI) * time.Second
	} else {
		cfg.ReportInterval = time.Duration(*rep) * time.Second
	}

	return cfg, nil
}