package v1

import (
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// GetMetric получить все митрики.
func (h V1) GetMetric(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetMetric(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, res)
	if err != nil {
		log.Err(err).Msg("fmt.Fprint")
	}
}

// GetOneMetric получить одну митрику.
func (h V1) GetOneMetric(w http.ResponseWriter, r *http.Request) {
	m := chi.URLParam(r, "metric")
	n := chi.URLParam(r, "name")

	if _, ok := service.TypeMetric[m]; !ok {
		http.Error(w, errors.New("invalid request data error").Error(), http.StatusBadRequest)
		return
	}

	if n == "" {
		http.Error(w, errors.New("invalid request value error").Error(), http.StatusBadRequest)
		return
	}

	req := service.MetricsRequest{
		ID:    n,
		MType: m,
	}

	res, err := h.service.GetOneMetric(r.Context(), req)
	if err != nil {
		if errors.Is(err, memory.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if m == service.Gauge {
		_, err = fmt.Fprint(w, *res.Value)
		log.Err(err).Msg("fmt.Fprint")
		return
	}

	_, err = fmt.Fprint(w, *res.Delta)
	if err != nil {
		log.Err(err).Msg("fmt.Fprint")
	}
}

// UpdateMetric загрузить одну митрику.
func (h V1) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	m := chi.URLParam(r, "metric")

	if _, ok := service.TypeMetric[m]; !ok {
		http.Error(w, errors.New("invalid request data error").Error(), http.StatusBadRequest)
		return
	}

	v := chi.URLParam(r, "value")
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		http.Error(w, errors.New("invalid request value error").Error(), http.StatusBadRequest)
		return
	}

	n := chi.URLParam(r, "name")

	intValue := int64(f)
	req := service.MetricsRequest{
		ID:    n,
		MType: m,
		Value: &f,
		Delta: &intValue,
	}
	if _, ok := service.TypeMetric[req.MType]; !ok {
		http.Error(w, errors.New("validate request type metric error").Error(), http.StatusBadRequest)
		return
	}

	if _, err = h.service.UpdateMetric(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, v)
	if err != nil {
		log.Err(err).Msg("fmt.Fprint")
	}
}

// GetOneMetric2 получить одну митрику.
func (h V1) GetOneMetric2(w http.ResponseWriter, r *http.Request) {
	l := h.handler.Log.With().Str("v1/metric", "GetOneMetric2").Logger()

	var req service.MetricsRequest
	var reader io.Reader
	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			l.Error().Err(err).Msgf("gzip.NewReader, request value: %+v", r)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() { err = gz.Close() }()

		reader = gz
	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		l.Error().Err(err).Msgf("io.ReadAll, request value: %+v", r)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := service.TypeMetric[req.MType]; !ok {
		err = errors.New("invalid request type metric error")
		l.Error().Err(err).Msgf("service.TypeMetric[metric], query param metric: %s", req.MType)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.GetOneMetric(r.Context(), req)
	if err != nil {
		l.Error().Err(err).Msgf("h.service.GetOneMetric, request value: %+v", req)

		if errors.Is(err, memory.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(res)
	if err != nil {
		l.Err(err).Msg("w.Write")
	}
}

// UpdateMetric2 загрузить батч митрик.
func (h V1) UpdateMetric2(w http.ResponseWriter, r *http.Request) {
	l := h.handler.Log.With().Str("v1/metric", "UpdateMetric2").Logger()

	var req service.MetricsRequest
	var reader io.Reader
	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			l.Error().Err(err).Msgf("gzip.NewReader, request value: %+v", r)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() { err = gz.Close() }()
		reader = gz
	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		l.Error().Err(err).Msgf("io.ReadAll, request value: %+v", r)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := service.TypeMetric[req.MType]; !ok {
		err = errors.New("invalid request type metric error")
		l.Error().Err(err).Msgf("service.TypeMetric[metric], query param metric: %s", req.MType)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.UpdateMetric(r.Context(), req)
	if err != nil {
		l.Error().Err(err).Msgf("h.service.UpdateMetric, request value: %+v", req)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(res)
	if err != nil {
		l.Err(err).Msg("w.Write")
	}
}

// Ping проверка подключения базы данных.
func (h V1) Ping(w http.ResponseWriter, r *http.Request) {
	l := h.handler.Log.With().Str("v1/metric", "Ping").Logger()

	err := h.service.Ping(r.Context())
	if err != nil {
		l.Error().Err(err).Msg("postgres.New")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(""))
	if err != nil {
		l.Err(err).Msg("w.Write")
	}
}

// Updates обновить метрику.
func (h V1) Updates(w http.ResponseWriter, r *http.Request) {
	l := h.handler.Log.With().Str("v1/metric", "UpdateMetric2").Logger()

	var reader io.Reader
	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			l.Error().Err(err).Msgf("gzip.NewReader, request value: %+v", r)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() { err = gz.Close() }()
		reader = gz
	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		l.Error().Err(err).Msgf("io.ReadAll, request value: %+v", r)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := os.ReadFile(h.cryptoKey)
	if err != nil {
		l.Error().Err(err).Msg("os.ReadFile")
		return
	}

	block, rest := pem.Decode(data)
	if block == nil {
		fmt.Println("Failed to decode PEM block")
		if rest != nil {
			fmt.Printf("Remaining data after PEM block: %s\n", rest)
		}
		return
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		l.Error().Err(err).Msg(" x509.ParsePKCS1PrivateKey")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	encrypteData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, body)
	if err != nil {
		l.Error().Err(err).Msg(" rsa.DecryptPKCS1v15")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var req []service.MetricsRequest
	err = json.Unmarshal(encrypteData, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Updates(r.Context(), req)
	if err != nil {
		l.Error().Err(err).Msgf("h.service.Updates, request value: %+v", req)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(""))
	if err != nil {
		l.Err(err).Msg("w.Write")
	}
}
