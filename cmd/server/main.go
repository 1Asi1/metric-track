package main

import (
	"flag"
	"fmt"

	"github.com/1Asi1/metric-track.git/internal/config"
	"github.com/1Asi1/metric-track.git/internal/server/apiserver"
)

var host string

func parseFlag() {
	flag.StringVar(&host, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
}
func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
	}

	parseFlag()
	cfg.MetricServerAddr = host
	server := apiserver.New(cfg)
	server.Run()
}
