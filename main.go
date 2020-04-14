package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

//Meta Data
var (
	Name    = "x-streaming"
	Address = ":8080"
)

func main() {

	srv := micro.NewService(micro.Version("v1.0.0"))

	srv.Init()

	go Run("config.toml")
	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
