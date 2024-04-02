package integration

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	random "math/rand"
	"net/http"
	"os"
	"time"

	"github.com/1Asi1/metric-track.git/internal/agent/config"
	"github.com/1Asi1/metric-track.git/internal/agent/service"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type MetricsRequest struct {
	MType string `json:"type"`
	Delta any    `json:"delta"`
	Value any    `json:"value"`
	ID    string `json:"id"`
}

type Client struct {
	service service.Service
	http    *resty.Client
	log     zerolog.Logger
	cfg     config.Config
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
	res = c.service.GetMetric()

	tickerPool := time.NewTicker(c.cfg.PollInterval)
	tickerRep := time.NewTicker(c.cfg.ReportInterval)

	type data map[string]any
	job := make(chan data, len(res.Type))
	counter := make(chan int, len(res.Type))
	for i := 0; i < c.cfg.RateLimit; i++ {
		go func(job chan data, counter chan int) {
			for j := range job {
				ct := <-counter
				if err := c.sendToServerBatch(j, ct); err != nil {
					l.Error().Err(err).Msgf("c.sendToServerBatch")
					return
				}
			}
		}(job, counter)
	}

	for {
		select {
		case <-tickerPool.C:
			res = c.service.GetMetric()
			count++
			res.Type["RandomValue"] = random.ExpFloat64()
		case <-tickerRep.C:
			for k, v := range res.Type {
				req := make(data, 1)
				req[k] = v
				job <- req
				counter <- count
			}
			count = 0
		}
	}
}

func (c *Client) sendToServerBatch(req map[string]any, count int) error {
	metrics := make([]MetricsRequest, 2)
	for k, v := range req {
		metrics[0] = MetricsRequest{
			ID:    k,
			MType: "gauge",
			Value: v,
			Delta: nil,
		}

		metrics[1] = MetricsRequest{
			ID:    k,
			MType: "counter",
			Delta: count,
			Value: nil,
		}
	}

	url := fmt.Sprintf("http://%s/updates/", c.cfg.MetricServerAddr)

	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(c.cfg.CryptoKey, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	key := make([]byte, 1024)
	_, err = file.Read(key)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(key)
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return err
	}
	encrypteData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	if err != nil {
		return err
	}
	defer func() { _ = gz.Close() }()
	_, err = gz.Write(encrypteData)
	if err != nil {
		return err
	}
	_ = gz.Close()

	request := c.http.R().SetHeader("Content-Type", "application/json")
	request.SetHeader("Content-Encoding", "gzip")

	request.SetBody(&buf)
	request.Method = resty.MethodPost
	request.URL = url
	defer c.http.SetCloseConnection(true)

	h1 := hmac.New(sha256.New, []byte(c.cfg.SecretKey))
	_, err = h1.Write(buf.Bytes())
	if err != nil {
		return err
	}
	res := hex.EncodeToString(h1.Sum(nil))
	request.SetHeader("HashSHA256", res)

	//зашифровать байты публичным ключем
	resp, err := request.Send()
	if err != nil {
		return err
	}
	defer func() { _ = resp.RawBody().Close() }()

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("expected status %d, got: %d", http.StatusOK, resp.StatusCode())
	}

	return nil
}
