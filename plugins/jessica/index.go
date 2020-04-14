package jessica

import (
	"log"
	"net/http"

	"github.com/micro-community/x-streaming/engine"
)

var config = new(engine.ListenerConfig)

func init() {
	engine.InstallPlugin(&engine.PluginConfig{
		Name:   "Jessica",
		Type:   engine.PLUGIN_SUBSCRIBER,
		Config: config,
		Run:    run,
	})
}
func run() {
	log.Printf("server Jessica start at %s", config.ListenAddr)
	log.Fatal(http.ListenAndServe(config.ListenAddr, http.HandlerFunc(WsHandler)))
}
