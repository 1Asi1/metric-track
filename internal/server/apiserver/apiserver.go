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
	Mux *chi.Mux
}

func New(cfg config.Config) APIServer {
	return APIServer{
		cfg: cfg,
		Mux: chi.NewRouter(),
	}
}

func (s *APIServer) Run() {
	memoryStore := memory.New("internal/server/repository/memory/store.json")
	metricS := service.New(memoryStore)
	route := rest.New(s.Mux, metricS)
	v1.New(route)

	log.Printf("server start: http://%s\n", s.cfg.MetricServerAddr)
	if err := http.ListenAndServe(s.cfg.MetricServerAddr, route.Mux); err != nil {
		log.Panicf("http.ListenAndServe panic: %v", err)
	}
}
