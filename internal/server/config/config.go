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
	StoreInterval    time.Duration
	StorePath        string
	StoreRestore     bool
}

func New(log zerolog.Logger) (Config, error) {
	l := log.With().Str("config", "New").Logger()

	var cfg Config
	add := flag.String("a", "localhost:8080", "address and port to run agent")
	store := flag.Int("i", 300, "store interval")
	path := flag.String("f", "./tmp/metrics-db.json", "path store file")
	restore := flag.Bool("r", true, "store restore")
	flag.Parse()

	metricServerAddrEnv, ok := os.LookupEnv("ADDRESS")
	if ok {
		l.Info().Msgf("server address value: %s", metricServerAddrEnv)
		cfg.MetricServerAddr = metricServerAddrEnv
	} else {
		l.Info().Msgf("server address value: %s", *add)
		cfg.MetricServerAddr = *add
	}

	storeIntervalEnv, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		sI, err := strconv.Atoi(storeIntervalEnv)
		if err != nil {
			l.Error().Err(err).Msgf("strconv.Atoi, store interval value: %s", storeIntervalEnv)
			return Config{}, err
		}

		cfg.StoreInterval = time.Duration(sI) * time.Second
	} else {
		cfg.StoreInterval = time.Duration(*store) * time.Second
	}

	storePathEnv, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		l.Info().Msgf("store path: %s", storePathEnv)
		cfg.StorePath = storePathEnv
	} else {
		l.Info().Msgf("store path: %s", *path)
		cfg.StorePath = *path
	}

	l.Info().Msgf("store restore: %v", *restore)
	cfg.StoreRestore = *restore

	return cfg, nil
}
