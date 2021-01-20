package rtsp

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/micro-community/stream/engine"
	"github.com/micro-community/stream/engine/avformat"
	"github.com/micro-community/stream/engine/util"
)

var collection = sync.Map{}
var config = struct {
	BufferLength int
	AutoPull     bool
	RemoteAddr   string
}{2048, true, "rtsp://localhost/${streamPath}"}

func init() {
	engine.InstallPlugin(&engine.PluginConfig{
		Name:    "RTSP",
		Type:    engine.PLUGIN_PUBLISHER | engine.PLUGIN_HOOK,
		Version: "1.0.0",
		Config:  &config,
		Run:     runPlugin,
		HotConfig: map[string]func(interface{}){
			"AutoPull": func(value interface{}) {
				config.AutoPull = value.(bool)
			},
		},
	})
}
func runPlugin() {
	engine.OnSubscribeHooks.AddHook(func(s *engine.Subscriber) {
		if config.AutoPull && s.Publisher == nil {
			new(RTSP).Publish(s.StreamPath, strings.Replace(config.RemoteAddr, "${streamPath}", s.StreamPath, -1))
		}
	})
	http.HandleFunc("/rtsp/list", func(w http.ResponseWriter, r *http.Request) {
		sse := util.NewSSE(w, r.Context())
		var err error
		for tick := time.NewTicker(time.Second); err == nil; <-tick.C {
			var info []*RTSPInfo
			collection.Range(func(key, value interface{}) bool {
				rtsp := value.(*RTSP)
				pinfo := &rtsp.RTSPInfo
				pinfo.BufferRate = len(rtsp.OutGoing) * 100 / config.BufferLength
				info = append(info, pinfo)
				return true
			})
			err = sse.WriteJSON(info)
		}
	})
	http.HandleFunc("/rtsp/pull", func(w http.ResponseWriter, r *http.Request) {
		targetURL := r.URL.Query().Get("target")
		streamPath := r.URL.Query().Get("streamPath")
		var err error
		if err == nil {
			new(RTSP).Publish(streamPath, targetURL)
			w.Write([]byte(`{"code":0}`))
		} else {
			w.Write([]byte(fmt.Sprintf(`{"code":1,"msg":"%s"}`, err.Error())))
		}
	})
}

type RTSP struct {
	engine.Publisher
	*Client
	RTSPInfo
}
type RTSPInfo struct {
	SyncCount  int64
	Header     *string
	BufferRate int
	StreamInfo *engine.StreamInfo
}

func (rtsp *RTSP) run() {
	fuBuffer := []byte{}
	iframeHead := []byte{0x17, 0x01, 0, 0, 0}
	pframeHead := []byte{0x27, 0x01, 0, 0, 0}
	spsHead := []byte{0xE1, 0, 0}
	ppsHead := []byte{0x01, 0, 0}
	nalLength := []byte{0, 0, 0, 0}

	avcsent := false
	aacsent := false
	handleNALU := func(nalType byte, payload []byte, ts int64) {
		rtsp.SyncCount++
		vl := len(payload)
		switch nalType {
		case avformat.NALU_IDR_Picture:
			if !avcsent {
				r := bytes.NewBuffer([]byte{})
				r.Write(avformat.RTMP_AVC_HEAD)
				util.BigEndian.PutUint16(spsHead[1:], uint16(len(rtsp.SPS)))
				r.Write(spsHead)
				r.Write(rtsp.SPS)
				util.BigEndian.PutUint16(ppsHead[1:], uint16(len(rtsp.PPS)))
				r.Write(ppsHead)
				r.Write(rtsp.PPS)
				rtsp.PushVideo(0, r.Bytes())
				avcsent = true
			}
			r := bytes.NewBuffer([]byte{})
			util.BigEndian.PutUint24(iframeHead[2:], 0)
			r.Write(iframeHead)
			util.BigEndian.PutUint32(nalLength, uint32(vl))
			r.Write(nalLength)
			r.Write(payload)
			rtsp.PushVideo(uint32(ts), r.Bytes())
		case avformat.NALU_Non_IDR_Picture:
			r := bytes.NewBuffer([]byte{})
			util.BigEndian.PutUint24(pframeHead[2:], 0)
			r.Write(pframeHead)
			util.BigEndian.PutUint32(nalLength, uint32(vl))
			r.Write(nalLength)
			r.Write(payload)
			rtsp.PushVideo(uint32(ts), r.Bytes())
		}
	}
	for {
		select {
		case <-rtsp.Done():
			return
		case data, ok := <-rtsp.OutGoing:
			if ok && data[0] == 36 {
				if data[1] == 0 {
					cc := data[4] & 0xF
					//rtp header
					rtphdr := 12 + cc*4

					//packet time
					ts := (int64(data[8]) << 24) + (int64(data[9]) << 16) + (int64(data[10]) << 8) + (int64(data[11]))

					//packet number
					//packno := (int64(data[6]) << 8) + int64(data[7])
					data = data[4+rtphdr:]
					nalType := data[0] & 0x1F

					if nalType >= 1 && nalType <= 23 {
						handleNALU(nalType, data, ts)
					} else if nalType == 28 {
						isStart := data[1]&0x80 != 0
						isEnd := data[1]&0x40 != 0
						nalType := data[1] & 0x1F
						//nri := (data[1]&0x60)>>5
						nal := data[0]&0xE0 | data[1]&0x1F
						if isStart {
							fuBuffer = []byte{0}
						}
						fuBuffer = append(fuBuffer, data[2:]...)
						if isEnd {
							fuBuffer[0] = nal
							handleNALU(nalType, fuBuffer, ts)
						}
					}

				} else if data[1] == 2 {
					// audio
					if !aacsent {
						rtsp.PushAudio(0, append([]byte{0xAF, 0x00}, rtsp.AudioSpecificConfig...))
					}
					cc := data[4] & 0xF
					rtphdr := 12 + cc*4
					payload := data[4+rtphdr:]
					auHeaderLen := (int16(payload[0]) << 8) + int16(payload[1])
					auHeaderLen = auHeaderLen >> 3
					auHeaderCount := int(auHeaderLen / 2)
					var auLenArray []int
					for iIndex := 0; iIndex < int(auHeaderCount); iIndex++ {
						auHeaderInfo := (int16(payload[2+2*iIndex]) << 8) + int16(payload[2+2*iIndex+1])
						auLen := auHeaderInfo >> 3
						auLenArray = append(auLenArray, int(auLen))
					}
					startOffset := 2 + 2*auHeaderCount
					for _, auLen := range auLenArray {
						endOffset := startOffset + auLen
						addHead := []byte{0xAF, 0x01}
						rtsp.PushAudio(0, append(addHead, payload[startOffset:endOffset]...))
					}
				}
			}
		}
	}
}

//Publish a rtsp stream
func (rtsp *RTSP) Publish(streamPath string, rtspUrl string) (result bool) {
	if result = rtsp.Publisher.Publish(streamPath); result {
		rtsp.Type = "RTSP"
		rtsp.Client = NewClient(config.BufferLength)
		rtsp.RTSPInfo.Header = &rtsp.Client.Header
		if status, message := rtsp.Client.Client(rtspUrl); !status {
			log.Println(message)
			return false
		}
		collection.Store(streamPath, rtsp)
		go rtsp.run()
	}
	return
}
