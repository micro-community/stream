package main

import (
	"log"
	"net/http"

	"github.com/micro-community/stream/ws"
	"go-micro.dev/v4/web"
)

// Meta Data
var (
	Name    = "x-streaming"
	Address = ":8080"
)

// this is still a toy code block
func main() {

	//support websocket directly,by go-micro
	srv := web.NewService(web.Name("stream"))

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
