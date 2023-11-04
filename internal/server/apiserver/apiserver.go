package apiserver

import (
	"log"
	"net/http"

	"github.com/1Asi1/metric-track.git/internal/config"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest/v1"
)

type APIServer struct {
	cfg config.Config
	Mux *http.ServeMux
}

func New(cfg config.Config) APIServer {
	return APIServer{
		cfg: cfg,
		Mux: http.NewServeMux(),
	}
}

func (s *APIServer) Run() {
	memoryStore := memory.New("internal/server/repository/memory/store.json")
	metricS := service.New(memoryStore)
	route := rest.New(s.Mux, metricS)
	v1.New(route)

	log.Println("server start: http://127.0.0.1:8080")
	if err := http.ListenAndServe(s.cfg.MetricServerPort, route.Mux); err != nil {
		log.Panicf("http.ListenAndServe panic: %v", err)
	}
}
