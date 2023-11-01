package rest

import (
	"net/http"

	"github.com/1Asi1/metric-track.git/internal/server/service"
)

type Handler struct {
	Mux     *http.ServeMux
	Service service.Service
}

func New(mux *http.ServeMux, service service.Service) Handler {
	return Handler{
		Mux:     mux,
		Service: service,
	}
}
