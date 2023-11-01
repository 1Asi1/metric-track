package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/1Asi1/metric-track.git/internal/server/service"
)

func (h V1) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.New("validate request method error").Error(), http.StatusBadRequest)
		return
	}

	uriData := strings.Split(r.URL.Path, "/")
	if len(uriData) == 0 || len(uriData) != 5 {
		http.Error(w, errors.New("invalid request data error").Error(), http.StatusNotFound)
		return
	}

	f, err := strconv.ParseFloat(uriData[4], 64)
	if err != nil {
		http.Error(w, errors.New("invalid request value error").Error(), http.StatusBadRequest)
		return
	}

	var value service.Type
	value.Gauge = f
	value.Counter = int64(f)
	req := service.Request{
		Metric: uriData[2],
		Name:   uriData[3],
		Type:   value,
	}

	if _, ok := service.TypeMetric[req.Metric]; !ok {
		http.Error(w, errors.New("validate request type metric error").Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateMetric(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, req)
}
