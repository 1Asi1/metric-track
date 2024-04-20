package config

import (
	_ "embed"
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type ConfigFile struct {
	MetricServerAddr string `json:"address"`
	StoreInterval    string `json:"store_interval"`
	StorePath        string `json:"store_file"`
	StoreRestore     bool   `json:"restore"`
	PostgresConnDSN  string `json:"database_dsn"`
	CryptoKey        string `json:"crypto_key"`
	TrustedSubnet    string `json:"trusted_subnet"`
}

type Config struct {
	MetricServerAddr string
	MetricPPROFAddr  string
	StoreInterval    time.Duration
	StorePath        string
	StoreRestore     bool
	PostgresConnDSN  string
	SecretKey        string
	CryptoKey        string
	TrustedSubnet    string
}

func New(log zerolog.Logger) (Config, error) {
	l := log.With().Str("config", "New").Logger()

	cfgFile := flag.String("c", "config.json", "config name")
	add := flag.String("a", "", "address and port to run agent")
	addPPROF := flag.String("p", "", "address and port to run pprof")
	store := flag.Int("i", 0, "store interval")
	path := flag.String("f", "", "path store file")
	restore := flag.Bool("r", false, "store restore")
	postgresql := flag.String("d", "", "dsn connecting to postgres")
	key := flag.String("k", "", "secret key for server")
	cryptoKey := flag.String("crypto-key", "", "crypto key for agent")
	trusted := flag.String("t", "", "trusted-subnet")
	flag.Parse()

	var cfgPathName string
	cfgFileEnv, ok := os.LookupEnv("CONFIG")
	if ok {
		l.Info().Msgf("config value: %s", cfgFileEnv)
		cfgPathName = cfgFileEnv
	} else {
		l.Info().Msgf("config address value: %s", *cfgFile)
		cfgPathName = *cfgFile
	}

	file, err := os.OpenFile("internal/server/config/"+cfgPathName, os.O_RDONLY, 0644)
	if err != nil {
		return Config{}, err
	}
	defer func() { _ = file.Close() }()

	var cfgFileData ConfigFile
	if err = json.NewDecoder(file).Decode(&cfgFileData); err != nil {
		return Config{}, err
	}

	var cfg Config
	metricServerAddrEnv, ok := os.LookupEnv("ADDRESS")
	if ok {
		l.Info().Msgf("server address value: %s", metricServerAddrEnv)
		cfg.MetricServerAddr = metricServerAddrEnv
	} else {
		l.Info().Msgf("server address value: %s", *add)
		cfg.MetricServerAddr = *add
		if cfg.MetricServerAddr == "" {
			cfg.MetricServerAddr = cfgFileData.MetricServerAddr
		}
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
		if cfg.StoreInterval == 0 {
			storeIntervalData, err := time.ParseDuration(cfgFileData.StoreInterval)
			if err != nil {
				return Config{}, err
			}

			cfg.StoreInterval = storeIntervalData
		}
	}

	storePathEnv, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		l.Info().Msgf("store path: %s", storePathEnv)
		cfg.StorePath = storePathEnv
	} else {
		l.Info().Msgf("store path: %s", *path)
		cfg.StorePath = *path
		if cfg.StorePath == "" {
			cfg.StorePath = cfgFileData.StorePath
		}
	}

	postgresqlAddrEnv, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		cfg.PostgresConnDSN = postgresqlAddrEnv
	} else {
		cfg.PostgresConnDSN = *postgresql
		if cfg.PostgresConnDSN == "" {
			cfg.PostgresConnDSN = cfgFileData.PostgresConnDSN
		}
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
		if cfg.CryptoKey == "" {
			cfg.CryptoKey = cfgFileData.CryptoKey
		}
	}

	trustedSubnetEnv, ok := os.LookupEnv("TRUSTED_SUBNET")
	if ok {
		cfg.TrustedSubnet = trustedSubnetEnv
	} else {
		cfg.TrustedSubnet = *trusted
		if cfg.TrustedSubnet == "" {
			cfg.TrustedSubnet = cfgFileData.TrustedSubnet
		}
	}

	l.Info().Msgf("store restore: %v", *restore)
	cfg.StoreRestore = *restore
	if !cfg.StoreRestore {
		cfg.StoreRestore = cfgFileData.StoreRestore
	}

	return cfg, nil
}
