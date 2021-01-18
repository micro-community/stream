package logrotate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/micro-community/streaming/engine"
)

var config = new(LogRotate)

type FileInfo struct {
	Name string
	Size int64
}
type LogRotate struct {
	Path        string
	Size        int64
	Days        int
	file        *os.File
	currentSize int64
	createTime  time.Time
	hours       float64
	splitFunc   func() bool
	formatter   string
}

func init() {
	engine.InstallPlugin(&engine.PluginConfig{
		Name:   "LogRotate",
		Type:   engine.PLUGIN_HOOK,
		Config: config,
		Run:    run,
	})
}
func run() {
	http.HandleFunc("/logrotate/tail", watchLogs)
	http.HandleFunc("/logrotate/find", findLog)
	http.HandleFunc("/logrotate/list", listLogFiles)
	http.HandleFunc("/logrotate/open", openLog)
	http.HandleFunc("/logrotate/download", download)
	if config.Size > 0 {
		config.splitFunc = config.splitBySize
	} else {
		if config.Days == 0 {
			config.Days = 1
		}
		config.hours = float64(config.Days) * 24
		config.splitFunc = config.splitByTime
	}
	config.createTime = time.Now()
	if runtime.GOOS == "windows" {
		config.formatter = "2006-01-02T15-04-05"
	} else {
		config.formatter = "2006-01-02T15:04:05"
	}
	err := os.MkdirAll(config.Path, 0666)
	config.file, err = os.OpenFile(filepath.Join(config.Path, fmt.Sprintf("%s.log", config.createTime.Format(config.formatter))), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err == nil {
		stat, _ := config.file.Stat()
		config.currentSize = stat.Size()
		AddWriter(config)
	} else {
		log.Println(err)
	}
}
func (l *LogRotate) splitBySize() bool {
	return l.currentSize >= l.Size
}
func (l *LogRotate) splitByTime() bool {
	return time.Since(l.createTime).Hours() > l.hours
}
func (l *LogRotate) Write(data []byte) (n int, err error) {
	n, err = l.file.Write(data)
	l.currentSize += int64(n)
	if err == nil {
		if l.splitFunc() {
			l.createTime = time.Now()
			if file, err := os.OpenFile(filepath.Join(l.Path, fmt.Sprintf("%s.log", l.createTime.Format(config.formatter))), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666); err == nil {
				l.file = file
				l.currentSize = 0
			}
		}
	}
	return
}
func watchLogs(w http.ResponseWriter, r *http.Request) {
	AddWriter(NewSSE(w, r.Context()))
	<-r.Context().Done()
}
func findLog(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("grep", fmt.Sprintf("\"%s\"", r.URL.Query().Get("query")), config.Path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	w.Write([]byte(out.String()))
}
func listLogFiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	dir, err := os.Open(config.Path)
	defer func() {
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}()
	if err != nil {
		return
	}
	var files []os.FileInfo
	files, err = dir.Readdir(0)
	if err != nil {
		return
	}
	var fileInfos []*FileInfo
	for _, info := range files {
		fileInfos = append(fileInfos, &FileInfo{
			info.Name(), info.Size(),
		})
	}
	var bytes []byte
	bytes, err = json.Marshal(fileInfos)
	if err != nil {
		return
	}
	w.Write(bytes)
}
func download(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	file, err := os.Open(filepath.Join(config.Path, filename))
	defer func() {
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}()
	w.Header().Add("Content-Disposition", "attachment; filename="+filename)
	if err != nil {
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		return
	}
}
func openLog(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	file, err := os.Open(filepath.Join(config.Path, filename))
	defer func() {
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}()
	if err != nil {
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		return
	}
}
