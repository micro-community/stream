package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Meta Data
var (
	Name    = "x-streaming"
	Address = ":8080"
)

// this is still a toy code block
func main() {

	//support websocket directly,by go-micro
	srv := echo.New()

	// Middleware
	srv.Use(middleware.Logger())
	srv.Use(middleware.Recover())

	// static files
	//srv.Handle("/", http.FileServer(http.Dir("html")))
	// Handle websocket connection
	//srv.HandleFunc("/ws", ws.HandleConn)

	//toy code ,will be changed.
	go Run("config.toml")

	// Run service
	if err := srv.Start(Address); err != nil {
		log.Fatal(err)
	}
}
