package gateway

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/micro-community/x-streaming/engine"
	"github.com/micro-community/x-streaming/engine/util"
)

var (
	config        = new(engine.ListenerConfig)
	startTime     = time.Now()
	dashboardPath string
)

func init() {
	_, currentFilePath, _, _ := runtime.Caller(0)
	dashboardPath = path.Join(path.Dir(currentFilePath), "../../dashboard/dist")
	log.Println(dashboardPath)
	engine.InstallPlugin(&engine.PluginConfig{
		Name:   "GateWay",
		Type:   engine.PLUGIN_HOOK,
		Config: config,
		Run:    run,
	})
}
func run() {
	http.HandleFunc("/api/sysInfo", sysInfo)
	http.HandleFunc("/api/stop", stopPublish)
	http.HandleFunc("/api/summary", summary)
	http.HandleFunc("/api/logs", watchLogs)
	http.HandleFunc("/api/config", getConfig)
	http.HandleFunc("/", website)
	log.Printf("server gateway start at %s", config.ListenAddr)
	log.Fatal(http.ListenAndServe(config.ListenAddr, nil))
}
func getConfig(w http.ResponseWriter, r *http.Request) {
	w.Write(engine.ConfigRaw)
}
func watchLogs(w http.ResponseWriter, r *http.Request) {
	engine.AddWriter(util.NewSSE(w, r.Context()))
	<-r.Context().Done()
}
func stopPublish(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if streamPath := r.URL.Query().Get("stream"); streamPath != "" {
		if b, ok := engine.AllRoom.Load(streamPath); ok {
			b.(*engine.Room).Cancel()
			w.Write([]byte("success"))
		} else {
			w.Write([]byte("no query stream"))
		}
	} else {
		w.Write([]byte("no such stream"))
	}
}
func website(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	} else if filePath == "/docs" {
		filePath = "/docs/index.html"
	}
	if mime := mime.TypeByExtension(path.Ext(filePath)); mime != "" {
		w.Header().Set("Content-Type", mime)
	}
	if f, err := ioutil.ReadFile(dashboardPath + filePath); err == nil {
		if _, err = w.Write(f); err != nil {
			w.WriteHeader(505)
		}
	} else {
		w.Header().Set("Location", "/")
		w.WriteHeader(302)
	}
}
func summary(w http.ResponseWriter, r *http.Request) {
	sse := util.NewSSE(w, r.Context())
	engine.Summary.Add()
	defer engine.Summary.Done()
	sse.WriteJSON(&engine.Summary)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if sse.WriteJSON(&engine.Summary) != nil {
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}
func sysInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	bytes, err := json.Marshal(engine.EngineInfo)
	if err == nil {
		_, err = w.Write(bytes)
	}
}
