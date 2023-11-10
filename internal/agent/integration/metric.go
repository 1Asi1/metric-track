package integration

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/1Asi1/metric-track.git/internal/agent/service"
	"github.com/1Asi1/metric-track.git/internal/config"
	s "github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	cfg     config.Config
	service service.Service
	http    *resty.Client
}

func New(cfg config.Config, s service.Service) *Client {
	return &Client{
		cfg:     cfg,
		service: s,
		http:    resty.New(),
	}
}

func (c *Client) SendMetricPeriodic() {
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
			if err := c.sendToServerGauge(res); err != nil {
				log.Println(err)
			}

			if err := c.sendToServerCounter(res); err != nil {
				log.Println(err)
			}

			count = 0
		}
	}
}

func (c *Client) sendToServerGauge(data service.Metric) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%f", c.cfg.MetricServerAddr, s.Gauge, "Name", data.RandomValue)

	if err := c.send(url); err != nil {
		return fmt.Errorf("sendToServerGauge: %v", err)
	}

	return nil
}

func (c *Client) sendToServerCounter(data service.Metric) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%d", c.cfg.MetricServerAddr, s.Counter, "Name", data.PollCount)

	if err := c.send(url); err != nil {
		return fmt.Errorf("sendToServerCounter: %v", err)
	}

	return nil
}

func (c *Client) send(url string) error {
	res, err := c.http.R().SetHeader("Content-Type", "text/plain; charset=utf-8").Post(url)
	if err != nil {
		return err
	}
	defer res.RawBody().Close()

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("error: %w status: %d", errors.New("status not ok"), res.StatusCode())
	}

	return nil
}
