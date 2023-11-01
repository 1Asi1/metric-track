package v1

import (
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
)

type V1 struct {
	handler rest.Handler
	Service service.Service
}

func New(h rest.Handler) {
	v1 := V1{
		handler: h,
		Service: h.Service,
	}

	v1.registerV1Route()
}

func (h V1) registerV1Route() {
	h.handler.Mux.HandleFunc("/update/", h.UpdateMetric)
}
