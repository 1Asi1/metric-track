package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/1Asi1/metric-track.git/internal/logger"
	"github.com/1Asi1/metric-track.git/internal/server/apiserver"
	"github.com/1Asi1/metric-track.git/internal/server/config"
	"github.com/rs/zerolog"
)

func main() {
	cfg, err := config.New(logger.NewLogger())
	if err != nil {
		log.Fatal("config.New")
	}

	l := logger.NewLogger()
	l = l.Level(zerolog.InfoLevel).With().Timestamp().Logger()

	go func() {
		if err = http.ListenAndServe(cfg.MetricPPROFAddr, nil); err != nil {
			l.Fatal().Err(err).Msg("http.ListenAndServe")
		}
	}()

	go func() {
		time.Sleep(20 * time.Second)
		fmem, err := os.Create("./profiles/mem.pprof")
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = fmem.Close()
		}()
		runtime.GC()
		if err = pprof.WriteHeapProfile(fmem); err != nil {
			panic(err)
		}
	}()

	server := apiserver.New(cfg, l)
	if err = server.Run(); err != nil {
		l.Fatal().Err(err).Msg("server.Run")
	}
}
