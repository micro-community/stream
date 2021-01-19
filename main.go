package main

import (
	"net/http"

	"github.com/micro-community/stream/ws"
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
	srv.Handle("/", http.FileServer(http.Dir("html")))

	// Handle websocket connection
	srv.HandleFunc("/ws", ws.HandleConn)

	//toy code ,will be changed.
	go Run("config.toml")

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
