package jessica

import (
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/micro-community/x-streaming/engine"
	"github.com/micro-community/x-streaming/engine/util"
)

var config = new(engine.ListenerConfig)
var publicPath string

func init() {
	publicPath = util.CurrentDir("dashboard", "public")
	engine.InstallPlugin(&engine.PluginConfig{
		Name:   "Jessica",
		Type:   engine.PLUGIN_SUBSCRIBER,
		Config: config,
		//		UI:      util.CurrentDir("dashboard", "ui", "plugin-jessica.min.js"),
		Version: "1.0.0",
		Run:     run,
	})
}
func run() {
	log.Printf("server Jessica start at %s", config.ListenAddr)
	http.HandleFunc("/jessibuca/", jessibuca)
	log.Fatal(http.ListenAndServe(config.ListenAddr, http.HandlerFunc(WsHandler)))
}
func jessibuca(w http.ResponseWriter, r *http.Request) {
	filePath := strings.TrimPrefix(r.URL.Path, "/jessibuca")
	if mime := mime.TypeByExtension(path.Ext(filePath)); mime != "" {
		w.Header().Set("Content-Type", mime)
	}
	if f, err := ioutil.ReadFile(publicPath + filePath); err == nil {
		if _, err = w.Write(f); err != nil {
			w.WriteHeader(500)
		}
	} else {
		w.WriteHeader(404)
	}
}
