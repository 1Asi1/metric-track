package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

const (
	intervalReport = 0
	intervalPull   = 0
	rateLimit      = 0
)

type ConfigFile struct {
	MetricServerAddr string `json:"address"`
	PollInterval     string `json:"poll_interval"`
	ReportInterval   string `json:"report_interval"`
	CryptoKey        string `json:"crypto_key"`
	GrpcAddr         string `json:"grpc_addr"`
}

type Config struct {
	MetricServerAddr string
	PollInterval     time.Duration
	ReportInterval   time.Duration
	SecretKey        string
	RateLimit        int
	CryptoKey        string
	ServerGrpcAddr   string
}

func New(log zerolog.Logger) (Config, error) {
	l := log.With().Str("config", "New").Logger()

	cfgFile := flag.String("c", "config.json", "config name")
	add := flag.String("a", "", "address and port to run agent")
	rep := flag.Int("r", intervalReport, "report agent interval")
	pull := flag.Int("p", intervalPull, "pull agent interval")
	key := flag.String("k", "", "secret key for agent")
	rLimit := flag.Int("l", rateLimit, "pull worker")
	cryptoKey := flag.String("crypto-key", "internal/agent/config/pbkey.pem", "crypto key for agent")
	grpc := flag.String("g", "127.0.0.1:8083", "grpc address")
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

	file, err := os.OpenFile("internal/agent/config/"+cfgPathName, os.O_RDONLY, 0644)
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
		pollIntervalData, err := time.ParseDuration(cfgFileData.PollInterval)
		if err != nil {
			return Config{}, err
		}
		if cfg.PollInterval == 0 {
			cfg.PollInterval = pollIntervalData
		}
	}

	reportIntervalEnv, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		rI, err := strconv.Atoi(reportIntervalEnv)
		if err != nil {
			return Config{}, fmt.Errorf("strconv.Atoi: %w", err)
		}

		cfg.ReportInterval = time.Duration(rI) * time.Second
	} else {
		reportIntervalData, err := time.ParseDuration(cfgFileData.ReportInterval)
		if err != nil {
			return Config{}, err
		}
		cfg.ReportInterval = time.Duration(*rep) * time.Second
		if cfg.ReportInterval == 0 {
			cfg.ReportInterval = reportIntervalData
		}
	}

	secretKeyEnv, ok := os.LookupEnv("KEY")
	if ok {
		cfg.SecretKey = secretKeyEnv
	} else {
		cfg.SecretKey = *key
	}

	rateLimitEnv, ok := os.LookupEnv("RATE_LIMIT")
	if ok {
		rL, err := strconv.Atoi(rateLimitEnv)
		if err != nil {
			return Config{}, fmt.Errorf("strconv.Atoi: %w", err)
		}

		cfg.RateLimit = rL
	} else {
		cfg.RateLimit = *rLimit
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

	grpcEnv, ok := os.LookupEnv("CONTENT_GRPC_ADDR")
	if ok {
		cfg.ServerGrpcAddr = grpcEnv
	} else {
		cfg.ServerGrpcAddr = *grpc
		if cfg.ServerGrpcAddr == "" {
			cfg.ServerGrpcAddr = cfgFileData.GrpcAddr
		}
	}

	return cfg, nil
}
