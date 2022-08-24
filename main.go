package main

import (
	"context"
	"net/http"

	"github.com/micro-community/stream/app"
	"github.com/micro-community/stream/util"
	"github.com/micro-community/stream/ws"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/util/log"
	"go-micro.dev/v4/web"
)

// Meta Data
var (
	Name    = "x-streaming"
	Address = ":8080"
)

// this is still a toy code block
func main() {

	//replace inner service
	s := micro.NewService(
		micro.Name("stream"),
		micro.HandleSignal(false),
		micro.Flags(
			&cli.StringFlag{
				Name: "version",
			},
			&cli.BoolFlag{
				Name: "timeout",
			},
		),
	)

	s.Init(micro.Action(func(c *cli.Context) error {
		app.Version = c.String("version")
		return nil
	}))

	//support websocket directly,by go-micro
	srv := web.NewService(web.MicroService(s))
	// static files
	srv.Handle("/", http.FileServer(http.Dir("html")))
	// Handle websocket connection
	srv.HandleFunc("/ws", ws.HandleConn)

	srv.Init()

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "version", app.Version))
	go util.WaitTerm(cancel)
	//toy code ,will be changed.
	go Run(ctx, "config.yaml")

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
