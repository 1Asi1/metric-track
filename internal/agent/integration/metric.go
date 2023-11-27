package integration

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/1Asi1/metric-track.git/internal/agent/service"
	"github.com/1Asi1/metric-track.git/internal/config"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type MetricsRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
	Delta any    `json:"delta"`
	Value any    `json:"value"`
}

type Client struct {
	cfg     config.Config
	service service.Service
	http    *resty.Client
	log     zerolog.Logger
}

func New(cfg config.Config, s service.Service, log zerolog.Logger) *Client {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	return &Client{
		cfg:     cfg,
		service: s,
		http:    client,
		log:     log,
	}
}

func (c *Client) SendMetricPeriodic() {
	l := c.log.With().Str("integration", "SendMetricPeriodic").Logger()

	var count int
	var res service.Metric
	tickerPool := time.NewTicker(c.cfg.PollInterval)
	tickerRep := time.NewTicker(c.cfg.ReportInterval)
	for {
		select {
		case <-tickerPool.C:
			res = c.service.GetMetric()

			count++

			res.Type["RandomValue"] = rand.ExpFloat64()
			res.Type["PollCount"] = count
		case <-tickerRep.C:

			for k, v := range res.Type {
				if err := c.sendToServerGauge(k, v); err != nil {
					l.Error().Err(err).Msgf("c.sendToServerGauge, type: %s, value: %v", k, v)
				}

				if err := c.sendToServerCounter(k, res.Type["PollCount"]); err != nil {
					l.Error().Err(err).Msgf("c.sendToServerGauge, type: %s, value: %v", k, v)
				}
			}

			count = 0
		}
	}
}

func (c *Client) sendToServerGauge(name string, value any) error {
	req := MetricsRequest{
		ID:    name,
		MType: "gauge",
		Value: value,
		Delta: 0,
	}

	if err := c.send(req); err != nil {
		return fmt.Errorf("sendToServerGauge: %v", err)
	}

	return nil
}

func (c *Client) sendToServerCounter(name string, value any) error {
	req := MetricsRequest{
		ID:    name,
		MType: "counter",
		Delta: value,
		Value: 0,
	}

	if err := c.send(req); err != nil {
		return fmt.Errorf("sendToServerCounter: %v", err)
	}

	return nil
}

func (c *Client) send(req MetricsRequest) error {
	l := c.log.With().Str("integration", "send").Logger()

	url := fmt.Sprintf("http://%s/update/", c.cfg.MetricServerAddr)

	data, err := json.Marshal(req)
	if err != nil {
		l.Error().Err(err).Msg("json.Marshal")
		return err
	}

	buf := bytes.NewBuffer(data)
	gz, err := gzip.NewWriterLevel(buf, gzip.BestSpeed)
	if err != nil {
		l.Error().Err(err).Msg("gzip.NewWriterLevel")
		return err
	}
	defer gz.Close()

	request := c.http.R().SetHeader("Content-Type", "application/json")
	request.SetHeader("Content-Encoding", "gzip")
	request.SetBody(data)
	request.Method = resty.MethodPost
	request.URL = url
	defer c.http.SetCloseConnection(true)

	res, err := request.Send()
	if err != nil {
		return err
	}
	defer res.RawBody().Close()

	if res.StatusCode() != http.StatusOK {
		l.Error().Err(err).Msgf("expected status %d, got: %d", http.StatusOK, res.StatusCode())
		return fmt.Errorf("expected status %d, got: %d", http.StatusOK, res.StatusCode())
	}

	return nil
}
