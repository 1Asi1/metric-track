package config

import (
	_ "embed"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	MetricServerAddr string
	MetricPPROFAddr  string
	StoreInterval    time.Duration
	StorePath        string
	StoreRestore     bool
	PostgresConnDSN  string
	SecretKey        string
	CryptoKey        string
}

func New(log zerolog.Logger) (Config, error) {
	l := log.With().Str("config", "New").Logger()

	var cfg Config
	add := flag.String("a", "localhost:8080", "address and port to run agent")
	addPPROF := flag.String("p", "localhost:8081", "address and port to run pprof")
	store := flag.Int("i", 300, "store interval")
	path := flag.String("f", "./tmp/metrics-db.json", "path store file")
	restore := flag.Bool("r", true, "store restore")
	postgresql := flag.String("d", "", "dsn connecting to postgres")
	key := flag.String("k", "", "secret key for server")
	cryptoKey := flag.String("crypto-key", "internal/server/config/ptkey.pem", "crypto key for agent")
	flag.Parse()

	metricServerAddrEnv, ok := os.LookupEnv("ADDRESS")
	if ok {
		l.Info().Msgf("server address value: %s", metricServerAddrEnv)
		cfg.MetricServerAddr = metricServerAddrEnv
	} else {
		l.Info().Msgf("server address value: %s", *add)
		cfg.MetricServerAddr = *add
	}

	metricPPROFAddrEnv, ok := os.LookupEnv("ADDRESS_PPROF")
	if ok {
		l.Info().Msgf("pprof address value: %s", metricServerAddrEnv)
		cfg.MetricPPROFAddr = metricPPROFAddrEnv
	} else {
		l.Info().Msgf("pprof address value: %s", *add)
		cfg.MetricPPROFAddr = *addPPROF
	}

	storeIntervalEnv, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		sI, err := strconv.Atoi(storeIntervalEnv)
		if err != nil {
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

	postgresqlAddrEnv, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		cfg.PostgresConnDSN = postgresqlAddrEnv
	} else {
		cfg.PostgresConnDSN = *postgresql
	}

	secretKeyEnv, ok := os.LookupEnv("KEY")
	if ok {
		cfg.SecretKey = secretKeyEnv
	} else {
		cfg.SecretKey = *key
	}

	cryptoKeyEnv, ok := os.LookupEnv("CRYPTO_KEY")
	if ok {
		cfg.CryptoKey = cryptoKeyEnv
	} else {
		cfg.CryptoKey = *cryptoKey
	}

	l.Info().Msgf("store restore: %v", *restore)
	cfg.StoreRestore = *restore

	return cfg, nil
}
