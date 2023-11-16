package rest

import (
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Handler struct {
	Mux     *chi.Mux
	Service service.Service
	Log     zerolog.Logger
}

func New(mux *chi.Mux, service service.Service, log zerolog.Logger) Handler {
	return Handler{
		Mux:     mux,
		Service: service,
		Log:     log,
	}
}
