package HDL

import (
	"log"
	"net/http"
	"strings"

	"github.com/micro-community/x-streaming/engine"
	"github.com/micro-community/x-streaming/engine/avformat"
)

var config = new(engine.ListenerConfig)

func init() {
	engine.InstallPlugin(&engine.PluginConfig{
		Name:   "HDL",
		Type:   engine.PLUGIN_SUBSCRIBER,
		Config: config,
		Run:    run,
	})
}

func run() {
	log.Printf("HDL start at %s", config.ListenAddr)
	log.Fatal(http.ListenAndServe(config.ListenAddr, http.HandlerFunc(HDLHandler)))
}

func HDLHandler(w http.ResponseWriter, r *http.Request) {
	sign := r.URL.Query().Get("sign")
	if err := engine.AuthHooks.Trigger(sign); err != nil {
		w.WriteHeader(403)
		return
	}
	stringPath := strings.TrimLeft(r.RequestURI, "/")
	if strings.HasSuffix(stringPath, ".flv") {
		stringPath = strings.TrimRight(stringPath, ".flv")
	}
	if _, ok := engine.AllRoom.Load(stringPath); ok {
		//atomic.AddInt32(&hdlId, 1)
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Content-Type", "video/x-flv")
		w.Write(avformat.FLVHeader)
		p := engine.Subscriber{
			Sign: sign,
			SendHandler: func(packet *avformat.SendPacket) error {
				return avformat.WriteFLVTag(w, packet)
			},
			SubscriberInfo: engine.SubscriberInfo{
				ID: r.RemoteAddr, Type: "FLV",
			},
		}
		p.Play(stringPath)
	} else {
		w.WriteHeader(404)
	}
}
