package integration

import (
	"fmt"
	"net/http"
	"time"

	"github.com/1Asi1/metric-track.git/internal/agent/service"
	"github.com/1Asi1/metric-track.git/internal/config"
	s "github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type Client struct {
	cfg     config.Config
	service service.Service
	http    *resty.Client
	log     zerolog.Logger
}

func New(cfg config.Config, s service.Service, log zerolog.Logger) *Client {
	return &Client{
		cfg:     cfg,
		service: s,
		http:    resty.New(),
		log:     log,
	}
}

func (c *Client) SendMetricPeriodic() {
	l := c.log.With().Str("integration", "SendMetricPeriodic").Logger()

	var count service.Counter
	var res service.Metric
	tickerPool := time.NewTicker(c.cfg.PollInterval)
	tickerRep := time.NewTicker(c.cfg.ReportInterval)
	for {
		select {
		case <-tickerPool.C:
			res = c.service.GetMetric()

			count++
			res.PollCount = count
		case <-tickerRep.C:

			for k, v := range res.Type {
				if err := c.sendToServerGauge(k, v); err != nil {
					l.Error().Err(err).Msgf("c.sendToServerGauge, type: %s, value: %v", k, v)
				}

				if err := c.sendToServerCounter(k, res.PollCount); err != nil {
					l.Error().Err(err).Msgf("c.sendToServerGauge, type: %s, value: %v", k, v)
				}
			}

			count = 0
		}
	}
}

func (c *Client) sendToServerGauge(name string, value any) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%v", c.cfg.MetricServerAddr, s.Gauge, name, value)

	if err := c.send(url); err != nil {
		return fmt.Errorf("sendToServerGauge: %v", err)
	}

	return nil
}

func (c *Client) sendToServerCounter(name string, value any) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%d", c.cfg.MetricServerAddr, s.Counter, name, value)

	if err := c.send(url); err != nil {
		return fmt.Errorf("sendToServerCounter: %v", err)
	}

	return nil
}

func (c *Client) send(url string) error {
	l := c.log.With().Str("integration", "send").Logger()

	res, err := c.http.R().SetHeader("Content-Type", "text/plain; charset=utf-8").Post(url)
	if err != nil {
		l.Error().Err(err).Msg("c.http.R()")
		return err
	}
	defer res.RawBody().Close()

	if res.StatusCode() != http.StatusOK {
		l.Error().Err(err).Msgf("expected status %d, got: %d", http.StatusOK, res.StatusCode())
		return fmt.Errorf("expected status %d, got: %d", http.StatusOK, res.StatusCode())
	}

	return nil
}
