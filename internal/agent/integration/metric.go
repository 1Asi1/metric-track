package integration

import (
	"bytes"
	"compress/gzip"
	"context"
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
	proto "github.com/1Asi1/metric-track.git/rpc/gen"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	grpc    proto.MetricGrpcClient
}

func New(cfg config.Config, s service.Service, log zerolog.Logger) *Client {
	client := resty.New()
	client.SetTimeout(10 * time.Second)

	contentConnDialCtx, contentConnDialCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer contentConnDialCancel()

	clientConn, err := grpc.DialContext(contentConnDialCtx, cfg.ServerGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Err(err).Msgf("grpc.Dial error: %v", err)
	}

	grpcClient := proto.NewMetricGrpcClient(clientConn)
	return &Client{
		cfg:     cfg,
		service: s,
		http:    client,
		log:     log,
		grpc:    grpcClient,
	}
}

func (c *Client) SendMetricPeriodic(ctx context.Context) {
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
				if err := c.sendToServerBatch(ctx, j, ct); err != nil {
					l.Error().Err(err).Msgf("c.sendToServerBatch")
					return
				}
				if err := c.sendToServerBatchGrpc(ctx, j, ct); err != nil {
					l.Error().Err(err).Msgf("c.sendToServerBatchGrpc")
					return
				}
			}
		}(job, counter)
	}

	go func() {
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
	}()
}

func (c *Client) sendToServerBatch(ctx context.Context, req map[string]any, count int) error {
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

	request := c.http.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Real-IP", c.http.BaseURL)
	request.SetContext(ctx)
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

func (c *Client) sendToServerBatchGrpc(ctx context.Context, req map[string]any, count int) error {
	metrics := make([]*proto.Metric, 2)
	for k, v := range req {
		metrics[0] = &proto.Metric{
			ID:    k,
			MType: "gauge",
			Value: v.(float64),
		}

		metrics[1] = &proto.Metric{
			ID:    k,
			MType: "counter",
			Delta: int64(count),
		}
	}

	_, err := c.grpc.Updates(ctx, &proto.UpdatesRequest{Metrics: []*proto.Metric{}}, nil)
	if err != nil {
		return err
	}

	return nil
}
