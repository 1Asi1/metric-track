package main

import (
	"fmt"

	"github.com/1Asi1/metric-track.git/internal/config"
	"github.com/1Asi1/metric-track.git/internal/server/apiserver"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
	}

	server := apiserver.New(cfg)
	server.Run()
}
