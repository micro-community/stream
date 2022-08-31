package webrtc

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/micro-community/stream/app"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v3"
)

var (
	reg_level = regexp.MustCompile("profile-level-id=(4.+f)")
)

type WebRTCOption struct {
	//config.Publish
	//config.Subscribe
	ICEServers []string
	PublicIP   []string
	PortMin    uint16
	PortMax    uint16
	PLI        time.Duration
	m          webrtc.MediaEngine
	s          webrtc.SettingEngine
	api        *webrtc.API
}

var webrtcOptions = &WebRTCOption{
	PLI: time.Second * 2,
}

// install pulgin
var webrtcObject = app.Install(webrtcOptions)

// TODO....
type FirstConfig struct{}

func (conf *WebRTCOption) OnEvent(event any) {
	switch event.(type) {
	case FirstConfig:
		RegisterCodecs(&conf.m)
		i := &interceptor.Registry{}
		if len(conf.PublicIP) > 0 {
			conf.s.SetNAT1To1IPs(conf.PublicIP, webrtc.ICECandidateTypeHost)
		}
		if conf.PortMin > 0 && conf.PortMax > 0 {
			conf.s.SetEphemeralUDPPortRange(conf.PortMin, conf.PortMax)
		}
		if len(conf.PublicIP) > 0 {
			conf.s.SetNAT1To1IPs(conf.PublicIP, webrtc.ICECandidateTypeHost)
		}
		conf.s.SetNetworkTypes([]webrtc.NetworkType{webrtc.NetworkTypeUDP4, webrtc.NetworkTypeUDP6})
		if err := webrtc.RegisterDefaultInterceptors(&conf.m, i); err != nil {
			panic(err)
		}
		conf.api = webrtc.NewAPI(webrtc.WithMediaEngine(&conf.m),
			webrtc.WithInterceptorRegistry(i), webrtc.WithSettingEngine(conf.s))
	}
}

func (conf *WebRTCOption) Play_(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/sdp")
	streamPath := r.URL.Path[len("/webrtc/play/"):]
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var suber = WebRTCSubscriber{
		WebRTCIO: WebRTCIO{SDP: string(bytes)},
	}
	if suber.PeerConnection, err = conf.api.NewPeerConnection(webrtc.Configuration{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	suber.OnICECandidate(func(ice *webrtc.ICECandidate) {
		if ice != nil {
			//	suber.Info(ice.ToJSON().Candidate)
		}
	})
	if err = suber.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: suber.SDP}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = webrtcObject.Subscribe(streamPath, &suber); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if sdp, err := suber.GetAnswer(); err == nil {
		w.Write([]byte(sdp))
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (conf *WebRTCOption) Push_(w http.ResponseWriter, r *http.Request) {
	streamPath := r.URL.Path[len("/webrtc/push/"):]
	w.Header().Set("Content-Type", "application/sdp")
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var puber = WebRTCPublisher{WebRTCIO: WebRTCIO{SDP: string(bytes)}}
	if puber.PeerConnection, err = conf.api.NewPeerConnection(webrtc.Configuration{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	puber.OnICECandidate(func(ice *webrtc.ICECandidate) {
		if ice != nil {
			//puber.Info(ice.ToJSON().Candidate)
		}
	})
	if _, err = puber.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = puber.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = webrtcObject.Publish(streamPath, &puber); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := puber.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: puber.SDP}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if answer, err := puber.GetAnswer(); err == nil {
		w.Write([]byte(answer))
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

var webrtcConfig = &WebRTCOption{
	PLI: time.Second * 2,
}
