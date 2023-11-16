package apiserver

import (
	"log"
	"net/http"

	"github.com/1Asi1/metric-track.git/internal/config"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest/v1"
	"github.com/go-chi/chi/v5"
)

type APIServer struct {
	cfg config.Config
	mux *chi.Mux
}

func New(cfg config.Config) APIServer {
	return APIServer{
		cfg: cfg,
		mux: chi.NewRouter(),
	}
}

func (s *APIServer) Run() error {
	memoryStore := memory.New()
	metricS := service.New(memoryStore)
	route := rest.New(s.mux, metricS)
	v1.New(route)

	log.Printf("server start: http://%s\n", s.cfg.MetricServerAddr)
	if err := http.ListenAndServe(s.cfg.MetricServerAddr, route.Mux); err != nil {
		return err
	}

	return nil
}
