package main

import "github.com/1Asi1/metric-track.git/internal/server/apiserver"

func main() {
	server := apiserver.New()
	server.Run()
}
