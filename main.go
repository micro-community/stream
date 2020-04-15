package main

import (
	"net/http"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
)

//Meta Data
var (
	Name    = "x-streaming"
	Address = ":8080"
)

// this is still a toy code block
func main() {

	//support websocket directly,by go-micro
	srv := web.NewService(web.Version("v1.0.0"))

	srv.Init()
	// static files
	srv.Handle("/websocket/", http.StripPrefix("/websocket/", http.FileServer(http.Dir("html"))))

	//toy code ,will be changed.
	go Run("config.toml")

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
