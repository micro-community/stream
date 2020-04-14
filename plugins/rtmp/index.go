package rtmp

import (
	"log"

	. "github.com/micro-community/x-streaming/engine"
)

var config = new(ListenerConfig)

func init() {
	InstallPlugin(&PluginConfig{
		Name:   "RTMP",
		Type:   PLUGIN_SUBSCRIBER | PLUGIN_PUBLISHER,
		Config: config,
		Run:    run,
	})
}
func run() {
	log.Printf("server rtmp start at %s", config.ListenAddr)
	log.Fatal(ListenRtmp(config.ListenAddr))
}
