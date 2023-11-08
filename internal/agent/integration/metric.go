package integration

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/1Asi1/metric-track.git/internal/agent/service"
	"github.com/1Asi1/metric-track.git/internal/config"
	s "github.com/1Asi1/metric-track.git/internal/server/service"
	"github.com/go-resty/resty/v2"
)

type Client interface {
	SendMetricPeriodic()
}

type client struct {
	cfg     config.Config
	service service.Service
	HTTP    *resty.Client
}

func New(cfg config.Config, s service.Service) Client {
	return client{
		cfg:     cfg,
		service: s,
		HTTP:    resty.New(),
	}
}

func (c client) SendMetricPeriodic() {
	var res service.Metric
	var count service.Counter
	var err error
	for i := 1; ; i++ {
		if i%c.cfg.PollInterval == 0 {
			res = c.service.GetMetric()

			count++
			res.PollCount = count
		}

		if i%c.cfg.ReportInterval == 0 {
			fmt.Printf("stat:  [value: %+v | count: %+v]\n", res.RandomValue, res.PollCount)

			if err = c.sendToServerGauge(res); err != nil {
				fmt.Println(err)
			}

			if err = c.sendToServerCounter(res); err != nil {
				fmt.Println(err)
			}

			count = 0
		}

		time.Sleep(1 * time.Second)
	}
}

func (c client) sendToServerGauge(data service.Metric) error {
	url := fmt.Sprintf("%s/update/%s/%s/%.6f", c.cfg.MetricServerAddr, s.Gauge, "Name", data.RandomValue)

	if err := c.send(url); err != nil {
		return fmt.Errorf("sendToServerGauge: %v", err)
	}

	return nil
}

func (c client) sendToServerCounter(data service.Metric) error {
	url := fmt.Sprintf("%s/update/%s/%s/%d", c.cfg.MetricServerAddr, s.Counter, "Name", data.PollCount)

	if err := c.send(url); err != nil {
		return fmt.Errorf("sendToServerCounter: %v", err)
	}

	return nil
}

func (c *client) send(url string) error {
	res, err := c.HTTP.R().SetHeader("Content-Type", "text/plain; charset=utf-8").Post(url)
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New("status not ok")
	}

	err = res.RawBody().Close()
	if err != nil {
		return err
	}

	return nil
}
