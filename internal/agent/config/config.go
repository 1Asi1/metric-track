package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

const (
	intervalReport = 10
	intervalPull   = 2
)

type Config struct {
	MetricServerAddr string
	PollInterval     time.Duration
	ReportInterval   time.Duration
	SecretKey        string
}

func New(log zerolog.Logger) (Config, error) {
	l := log.With().Str("config", "New").Logger()

	var cfg Config

	add := flag.String("a", "localhost:8080", "address and port to run agent")
	rep := flag.Int("r", intervalReport, "report agent interval")
	pull := flag.Int("p", intervalPull, "pull agent interval")
	key := flag.String("k", "", "secret key for agent")
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
			return Config{}, fmt.Errorf("strconv.Atoi: %w", err)
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
			return Config{}, fmt.Errorf("strconv.Atoi: %w", err)
		}

		cfg.ReportInterval = time.Duration(rI) * time.Second
	} else {
		cfg.ReportInterval = time.Duration(*rep) * time.Second
	}

	secretKeyEnv, ok := os.LookupEnv("KEY")
	if ok {
		cfg.SecretKey = secretKeyEnv
	} else {
		cfg.SecretKey = *key
	}

	return cfg, nil
}
