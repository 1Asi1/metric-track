package apiserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/1Asi1/metric-track.git/internal/server/config"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/1Asi1/metric-track.git/internal/server/repository/storage"
	"github.com/1Asi1/metric-track.git/internal/server/service"
	metric_grpc "github.com/1Asi1/metric-track.git/internal/server/transport/grpc"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest"
	"github.com/1Asi1/metric-track.git/internal/server/transport/rest/v1"
	proto "github.com/1Asi1/metric-track.git/rpc/gen"
	"github.com/go-chi/chi/v5"
	midlog "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

const (
	timeoutServerShutdown = time.Second * 5
)

type APIServer struct {
	cfg config.Config
	mux *chi.Mux
	log zerolog.Logger
}

func New(cfg config.Config, log zerolog.Logger) APIServer {
	return APIServer{
		cfg: cfg,
		mux: chi.NewRouter(),
		log: log,
	}
}

func (s *APIServer) Run() error {
	l := s.log.With().Str("apiserver", "Run").Logger()
	var store service.Store
	if s.cfg.PostgresConnDSN != "" {
		psqlCfg := storage.Config{
			ConnDSN:         s.cfg.PostgresConnDSN,
			Logger:          s.log,
			MaxConn:         30,
			MaxConnLifeTime: 10,
			MaxConnIdleTime: 10,
		}
		postgresql, err := storage.New(psqlCfg, s.log)
		if err != nil {
			l.Info().Msg("postgres.New error")
		}
		defer func() {
			if err = postgresql.Close(); err != nil {
				l.Err(err).Msg("postgresql.Close")
			}
		}()

		store = postgresql
	} else {
		store = memory.New(s.log, s.cfg)
	}

	metricS := service.New(store, s.log)
	route := rest.New(s.mux, metricS, s.log)

	route.Mux.Use(midlog.Logger)
	v1.New(route, s.cfg.SecretKey, s.cfg.CryptoKey)

	var srv = http.Server{Addr: s.cfg.MetricServerAddr}
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer l.Info().Msg("server has been shutdown")
		defer wg.Done()
		<-sigint
		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			l.Err(err).Msg("an error occurred during server shutdown")
		}
	}()

	go func() {
		grpcConn, err := net.Listen("tcp", fmt.Sprintf(":%v", s.cfg.GrpcPort))
		if err != nil {
			l.Err(err).Msgf("net.Listen error: %v; GrpcPort: %v", err, s.cfg.GrpcPort)
		}

		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(metric_grpc.CheckSubnetInterceptor(s.cfg.TrustedSubnet)),
			grpc.UnaryInterceptor(metric_grpc.HMACInterceptor(s.cfg.SecretKey)),
		)
		proto.RegisterMetricGrpcServer(grpcServer, metric_grpc.NewMetricGrpcServer(metricS))

		if err = grpcServer.Serve(grpcConn); err != nil {
			l.Err(err).Msgf("grpcServer.Serve error: %v; GrpcPort: %v", err, s.cfg.GrpcPort)
		}
	}()

	l.Info().Msgf("server start: http://%s", s.cfg.MetricServerAddr)
	srv.Handler = route.Mux
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		l.Error().Err(err).Msg("http.ListenAndServe")
		return err
	}
	return nil
}
