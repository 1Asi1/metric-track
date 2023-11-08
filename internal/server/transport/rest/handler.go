package rest

import (
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Mux     *chi.Mux
	Service service.Service
}

func New(mux *chi.Mux, service service.Service) Handler {
	return Handler{
		Mux:     mux,
		Service: service,
	}
}
