package v1

import (
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest/middleware"
	"github.com/go-chi/chi/v5"
)

type V1 struct {
	handler rest.Handler
	service service.Service
}

func New(h rest.Handler) {
	v1 := V1{
		handler: h,
		service: h.Service,
	}

	v1.handler.Mux.Use(middleware.GzipMiddleware)
	v1.registerV1Route()
}

func (h V1) registerV1Route() {
	h.handler.Mux.Route("/", func(r chi.Router) {
		r.Get("/", h.GetMetric)
		r.Get("/value/{metric}/{name}", h.GetOneMetric)
		r.Post("/update/{metric}/{name}/{value}", h.UpdateMetric)
		r.Post("/value/", h.GetOneMetric2)
		r.Post("/update/", h.UpdateMetric2)
	})
}
