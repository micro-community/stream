package cluster

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/micro-community/x-streaming/engine"
)

const (
	_ byte = iota
	MSG_AUDIO
	MSG_VIDEO
	MSG_SUBSCRIBE
	MSG_AUTH
	MSG_SUMMARY
	MSG_LOG
)

var (
	config = struct {
		Master     string
		ListenAddr string
	}{}
	slaves     = sync.Map{}
	masterConn *net.TCPConn
)

func init() {
	engine.InstallPlugin(&engine.PluginConfig{
		Name:   "Cluster",
		Type:   engine.PLUGIN_HOOK | engine.PLUGIN_PUBLISHER | engine.PLUGIN_SUBSCRIBER,
		Config: &config,
		Run:    run,
	})
}
func run() {
	if config.Master != "" {
		engine.OnSubscribeHooks.AddHook(onSubscribe)
		addr, err := net.ResolveTCPAddr("tcp", config.Master)
		if engine.MayBeError(err) {
			return
		}
		go readMaster(addr)
	}
	if config.ListenAddr != "" {
		engine.Summary.Children = make(map[string]*engine.ServerSummary)
		engine.OnSummaryHooks.AddHook(onSummary)
		log.Printf("server bare start at %s", config.ListenAddr)
		log.Fatal(ListenBare(config.ListenAddr))
	}
}

func readMaster(addr *net.TCPAddr) {
	var err error
	var cmd byte
	for {
		if masterConn, err = net.DialTCP("tcp", nil, addr); !engine.MayBeError(err) {
			reader := bufio.NewReader(masterConn)
			log.Printf("connect to master %s reporting", config.Master)
			for report(); err == nil; {
				if cmd, err = reader.ReadByte(); !engine.MayBeError(err) {
					switch cmd {
					case MSG_SUMMARY: //收到主服务器指令，进行采集和上报
						log.Println("receive summary request from master")
						if cmd, err = reader.ReadByte(); !engine.MayBeError(err) {
							if cmd == 1 {
								engine.Summary.Add()
								go onReport()
							} else {
								engine.Summary.Done()
							}
						}
					}
				}
			}
		}
		t := 5 + rand.Int63n(5)
		log.Printf("reconnect to master %s after %d seconds", config.Master, t)
		time.Sleep(time.Duration(t) * time.Second)
	}
}
func report() {
	if b, err := json.Marshal(engine.Summary); err == nil {
		data := make([]byte, len(b)+2)
		data[0] = MSG_SUMMARY
		copy(data[1:], b)
		data[len(data)-1] = 0
		_, err = masterConn.Write(data)
	}
}

//定时上报
func onReport() {
	for range time.NewTicker(time.Second).C {
		if engine.Summary.Running() {
			report()
		} else {
			return
		}
	}
}
func orderReport(conn io.Writer, start bool) {
	b := []byte{MSG_SUMMARY, 0}
	if start {
		b[1] = 1
	}
	conn.Write(b)
}

//通知从服务器需要上报或者关闭上报
func onSummary(start bool) {
	slaves.Range(func(k, v interface{}) bool {
		orderReport(v.(*net.TCPConn), start)
		return true
	})
}

func onSubscribe(s *engine.Subscriber) {
	if s.Publisher == nil {
		go PullUpStream(s.StreamPath)
	}
}
