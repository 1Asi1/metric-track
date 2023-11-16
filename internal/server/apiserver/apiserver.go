package apiserver

import (
	"net/http"

	"github.com/1Asi1/metric-track.git/internal/config"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest/v1"
	"github.com/go-chi/chi/v5"
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

	memoryStore := memory.New(s.log)
	metricS := service.New(memoryStore, s.log)
	route := rest.New(s.mux, metricS, s.log)
	v1.New(route)

	l.Info().Msgf("server start: http://%s", s.cfg.MetricServerAddr)
	if err := http.ListenAndServe(s.cfg.MetricServerAddr, route.Mux); err != nil {
		l.Error().Err(err).Msg("http.ListenAndServe")
		return err
	}

	return nil
}
