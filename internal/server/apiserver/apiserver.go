package apiserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/1Asi1/metric-track.git/internal/server/config"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/1Asi1/metric-track.git/internal/server/repository/storage"
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest/v1"
	"github.com/go-chi/chi/v5"
	midlog "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type APIServer struct {
	cfg config.Config
	mux *chi.Mux
	log zerolog.Logger
}

func New(cfg config.Config, log zerolog.Logger) APIServer {
	return APIServer{
		cfg: cfg,
		mux: chi.NewRouter(),
		log: log,
	}
}

func (s *APIServer) Run() error {
	l := s.log.With().Str("apiserver", "Run").Logger()
	var store service.Store
	if s.cfg.PostgresConnDSN != "" {
		psqlCfg := storage.Config{
			ConnDSN:         s.cfg.PostgresConnDSN,
			Logger:          s.log,
			MaxConn:         30,
			MaxConnLifeTime: 10,
			MaxConnIdleTime: 10,
		}
		postgresql, err := storage.New(psqlCfg, s.log)
		if err != nil {
			l.Info().Msg("postgres.New error")
		}
		defer func() {
			if err = postgresql.Close(); err != nil {
				l.Err(err).Msg("postgresql.Close")
			}
		}()

		store = postgresql
	} else {
		store = memory.New(s.log, s.cfg)
	}

	metricS := service.New(store, s.log)
	route := rest.New(s.mux, metricS, s.log)

	route.Mux.Use(midlog.Logger)
	v1.New(route, s.cfg.SecretKey, s.cfg.CryptoKey)

	var srv = http.Server{Addr: s.cfg.MetricServerAddr}
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			l.Error().Err(err).Msgf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	l.Info().Msgf("server start: http://%s", s.cfg.MetricServerAddr)
	srv.Handler = route.Mux
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		l.Error().Err(err).Msg("http.ListenAndServe")
		return err
	}
	<-idleConnsClosed
	l.Info().Msg("Start server shutdown gracefully...")
	return nil
}
