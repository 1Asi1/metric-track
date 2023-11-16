package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/go-chi/chi/v5"
)

func (h V1) GetMetric(w http.ResponseWriter, r *http.Request) {
	l := h.handler.Log.With().Str("v1/metric", "GetMetric").Logger()

	res, err := h.service.GetMetric(r.Context())
	if err != nil {
		l.Error().Err(err).Msg("h.service.GetMetric")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, res)
}

func (h V1) GetOneMetric(w http.ResponseWriter, r *http.Request) {
	l := h.handler.Log.With().Str("v1/metric", "GetOneMetric").Logger()

	m := chi.URLParam(r, "metric")
	n := chi.URLParam(r, "name")

	if _, ok := service.TypeMetric[m]; !ok {
		err := errors.New("invalid request type metric error")
		l.Error().Err(err).Msgf("service.TypeMetric[metric], query param metric: %s", m)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if n == "" {
		err := errors.New("invalid request name error")
		l.Error().Err(err).Msg("query param name is empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.service.GetOneMetric(r.Context(), m, n)
	if err != nil {
		l.Error().Err(err).Msgf("h.service.GetOneMetric, metric: %s, name: %s", m, n)

		if errors.Is(err, memory.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, res)
}

func (h V1) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	l := h.handler.Log.With().Str("v1/metric", "UpdateMetric").Logger()

	m := chi.URLParam(r, "metric")

	if _, ok := service.TypeMetric[m]; !ok {
		err := errors.New("invalid request type metric error")
		l.Error().Err(err).Msgf("service.TypeMetric[metric], query param metric: %s", m)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := chi.URLParam(r, "value")
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		l.Error().Err(err).Msgf("strconv.ParseFloat, query param value: %s", v)
		http.Error(w, errors.New("invalid request value error").Error(), http.StatusBadRequest)
		return
	}

	var value service.Type
	value.Gauge = f
	value.Counter = int64(f)
	n := chi.URLParam(r, "name")
	req := service.Request{
		Metric: m,
		Name:   n,
		Type:   value,
	}

	if _, ok := service.TypeMetric[req.Metric]; !ok {
		err = errors.New("invalid request type metric error")
		l.Error().Err(err).Msgf("service.TypeMetric[metric], query param metric: %s", m)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.service.UpdateMetric(r.Context(), req); err != nil {
		l.Error().Err(err).Msgf("h.service.UpdateMetric, req: %+v", req)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, v)
}
