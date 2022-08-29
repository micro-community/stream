package webrtc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/micro-community/stream/app"
	"github.com/micro-community/stream/engine"
	"github.com/micro-community/stream/engine/avformat"
	"github.com/micro-community/stream/service/rtp"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

var config struct {
	ICEServers []string
	PublicIP   []string
	PortMin    uint16
	PortMax    uint16
}

// var playWaitList WaitList
var reg_level = regexp.MustCompile("profile-level-id=(4.+f)")

type WaitList struct {
	m map[string]*WebRTC
	l sync.Mutex
}

func (wl *WaitList) Set(k string, v *WebRTC) {
	wl.l.Lock()
	defer wl.l.Unlock()
	if wl.m == nil {
		wl.m = make(map[string]*WebRTC)
	}
	wl.m[k] = v
}
func (wl *WaitList) Get(k string) *WebRTC {
	wl.l.Lock()
	defer wl.l.Unlock()
	defer delete(wl.m, k)
	return wl.m[k]
}
func init() {
	app.InstallPlugin(app.PluginOptions{
		Name: "WebRTC",
	})
}

type WebRTC struct {
	rtp.RTP
	*webrtc.PeerConnection
	RemoteAddr string
	audioTrack *webrtc.Track
	videoTrack *webrtc.Track
	m          webrtc.MediaEngine
	s          webrtc.SettingEngine
	api        *webrtc.API
	payloader  avformat.H264
	// codecs.H264Packet
	// *os.File
}

func (rtc *WebRTC) Play(streamPath string) bool {
	var sub engine.Subscriber
	sub.ID = rtc.RemoteAddr
	sub.Type = "WebRTC"
	var lastTimeStampV, lastTiimeStampA uint32
	sub.OnData = func(packet *avformat.SendPacket) error {
		if packet.Type == avformat.FLV_TAG_TYPE_AUDIO {
			var s uint32
			if lastTiimeStampA > 0 {
				s = packet.Timestamp - lastTiimeStampA
			}
			lastTiimeStampA = packet.Timestamp
			_ = rtc.audioTrack.WriteSample(media.Sample{
				Data:    packet.Payload[1:],
				Samples: s * 8,
			})
			return nil
		}
		if packet.IsSequence {
			rtc.payloader.PPS = sub.PPS
			rtc.payloader.SPS = sub.SPS
		} else {
			var s uint32
			if lastTimeStampV > 0 {
				s = packet.Timestamp - lastTimeStampV
			}
			lastTimeStampV = packet.Timestamp
			_ = rtc.videoTrack.WriteSample(media.Sample{
				Data:    packet.Payload,
				Samples: s * 90,
			})

		}
		return nil
	}
	// go sub.Subscribe(streamPath)
	rtc.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		engine.Printf("%s Connection State has changed %s ", streamPath, connectionState.String())
		switch connectionState {
		case webrtc.ICEConnectionStateDisconnected:
			sub.Close()
			rtc.Close()
		case webrtc.ICEConnectionStateConnected:

			//rtc.videoTrack = rtc.GetSenders()[0].Track()
			_ = sub.Subscribe(streamPath)
		}
	})
	return true
}
func (rtc *WebRTC) Publish(streamPath string) bool {
	rtc.m.RegisterCodec(webrtc.NewRTPCodec(webrtc.RTPCodecTypeVideo,
		webrtc.H264,
		90000,
		0,
		"level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42001f",
		webrtc.DefaultPayloadTypeH264,
		new(avformat.H264)))
	//m.RegisterCodec(NewRTPPCMUCodec(DefaultPayloadTypePCMU, 8000))
	if !strings.HasPrefix(rtc.RemoteAddr, "127.0.0.1") && !strings.HasPrefix(rtc.RemoteAddr, "[::1]") {
		rtc.s.SetNAT1To1IPs(config.PublicIP, webrtc.ICECandidateTypeHost)
	}
	if config.PortMin > 0 && config.PortMax > 0 {
		_ = rtc.s.SetEphemeralUDPPortRange(config.PortMin, config.PortMax)
	}
	rtc.api = webrtc.NewAPI(webrtc.WithMediaEngine(rtc.m), webrtc.WithSettingEngine(rtc.s))
	peerConnection, err := rtc.api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: config.ICEServers,
			},
		},
	})
	if err != nil {
		engine.Println(err)
		return false
	}
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		if err != nil {
			engine.Println(err)
			return false
		}
	}
	if err != nil {
		return false
	}
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		engine.Printf("%s Connection State has changed %s ", streamPath, connectionState.String())
		switch connectionState {
		case webrtc.ICEConnectionStateDisconnected, webrtc.ICEConnectionStateFailed:
			if rtc.Stream != nil {
				rtc.Stream.Close()
			}
		}
	})
	rtc.PeerConnection = peerConnection
	if rtc.RTP.Publish(streamPath) {
		//f, _ := os.OpenFile("resource/live/rtc.h264", os.O_TRUNC|os.O_WRONLY, 0666)
		rtc.Stream.Type = "WebRTC"
		peerConnection.OnTrack(func(track *webrtc.Track, receiver *webrtc.RTPReceiver) {
			defer rtc.Stream.Close()
			go func() {
				ticker := time.NewTicker(time.Second * 2)
				select {
				case <-ticker.C:
					if rtcpErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: track.SSRC()}}); rtcpErr != nil {
						engine.Println(rtcpErr)
					}
				case <-rtc.Done():
					return
				}
			}()
			pack := &rtp.RTPPack{
				Type: rtp.RTPType(track.Kind() - 1),
			}
			for b := make([]byte, 1460); ; rtc.PushPack(pack) {
				i, err := track.Read(b)
				if err != nil {
					return
				}
				if err = pack.Unmarshal(b[:i]); err != nil {
					return
				}
				// rtc.Unmarshal(pack.Payload)
				// f.Write(bytes)
			}
		})
	} else {
		return false
	}
	return true
}
func (rtc *WebRTC) GetAnswer() ([]byte, error) {
	// Sets the LocalDescription, and starts our UDP listeners
	answer, err := rtc.CreateAnswer(nil)
	if err != nil {
		return nil, err
	}
	if err := rtc.SetLocalDescription(answer); err != nil {
		engine.Println(err)
		return nil, err
	}
	if bytes, err := json.Marshal(answer); err != nil {
		engine.Println(err)
		return bytes, err
	} else {
		return bytes, nil
	}
}

func run() {
	http.HandleFunc("/webrtc/play", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		origin := r.Header["Origin"]
		if len(origin) == 0 {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", origin[0])
		}

		w.Header().Set("Content-Type", "application/json")
		streamPath := r.URL.Query().Get("streamPath")
		var offer webrtc.SessionDescription
		var rtc WebRTC

		bytes, err := ioutil.ReadAll(r.Body)
		defer func() {
			if err != nil {
				engine.Println(err)
				fmt.Fprintf(w, `{"errmsg":"%s"}`, err)
				return
			}
			rtc.Play(streamPath)
		}()
		if err != nil {
			return
		}
		if err = json.Unmarshal(bytes, &offer); err != nil {
			return
		}

		pli := "42001f"
		if stream := engine.FindStream(streamPath); stream != nil {
			<-stream.WaitPub
			pli = fmt.Sprintf("%x", stream.SPS[1:4])
		}
		if !strings.Contains(offer.SDP, pli) {
			pli = reg_level.FindAllStringSubmatch(offer.SDP, -1)[0][1]
		}
		rtc.m.RegisterCodec(webrtc.NewRTPCodec(webrtc.RTPCodecTypeVideo,
			webrtc.H264,
			90000,
			0,
			"level-asymmetry-allowed=1;packetization-mode=1;profile-level-id="+pli,
			webrtc.DefaultPayloadTypeH264,
			&rtc.payloader))
		rtc.m.RegisterCodec(webrtc.NewRTPPCMACodec(webrtc.DefaultPayloadTypePCMA, 8000))
		if !strings.HasPrefix(r.RemoteAddr, "127.0.0.1") && !strings.HasPrefix(r.RemoteAddr, "[::1]") {
			rtc.s.SetNAT1To1IPs(config.PublicIP, webrtc.ICECandidateTypeHost)
		}
		if config.PortMin > 0 && config.PortMax > 0 {
			_ = rtc.s.SetEphemeralUDPPortRange(config.PortMin, config.PortMax)
		}
		rtc.api = webrtc.NewAPI(webrtc.WithMediaEngine(rtc.m), webrtc.WithSettingEngine(rtc.s))

		if rtc.PeerConnection, err = rtc.api.NewPeerConnection(webrtc.Configuration{
			// ICEServers: []ICEServer{
			// 	{
			// 		URLs: config.ICEServers,
			// 	},
			// },
		}); err != nil {
			return
		}
		rtc.OnICECandidate(func(ice *webrtc.ICECandidate) {
			if ice != nil {
				engine.Println(ice.ToJSON().Candidate)
			}
		})

		rtc.RemoteAddr = r.RemoteAddr
		if err = rtc.SetRemoteDescription(offer); err != nil {
			return
		}

		if rtc.videoTrack, err = rtc.NewTrack(webrtc.DefaultPayloadTypeH264, 8, "video", "monibuca"); err != nil {
			return
		}
		if rtc.audioTrack, err = rtc.NewTrack(webrtc.DefaultPayloadTypePCMA, 9, "audio", "monibuca"); err != nil {
			return
		}
		if _, err = rtc.AddTrack(rtc.videoTrack); err != nil {
			return
		}
		if bytes, err := rtc.GetAnswer(); err == nil {
			_, _ = w.Write(bytes)
		} else {
			return
		}
	})

	http.HandleFunc("/webrtc/publish", func(w http.ResponseWriter, r *http.Request) {
		streamPath := r.URL.Query().Get("streamPath")
		offer := webrtc.SessionDescription{}
		bytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
		}
		err = json.Unmarshal(bytes, &offer)
		if err != nil {
			engine.Println(err)
			return
		}
		rtc := new(WebRTC)
		rtc.RemoteAddr = r.RemoteAddr
		if rtc.Publish(streamPath) {
			if err := rtc.SetRemoteDescription(offer); err != nil {
				engine.Println(err)
				return
			}
			if bytes, err = rtc.GetAnswer(); err == nil {
				_, _ = w.Write(bytes)
			} else {
				engine.Println(err)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		} else {
			_, _ = w.Write([]byte("bad name"))
		}
	})
}
