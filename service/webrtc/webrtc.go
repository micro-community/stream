package webrtc

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/micro-community/stream/app"
	"github.com/pion/interceptor"
	webrtc3 "github.com/pion/webrtc/v3"
)

var (
	reg_level = regexp.MustCompile("profile-level-id=(4.+f)")
	// init web rtc
	webrtcOptions = app.WebRTCOption{
		PLI:  time.Second * 2,
		Name: "WebRTC",
	}
)

// install webRTCPulgins
var webrtcObject = app.Install(app.WithWebRTC(webrtcOptions))

// webRTC
type webrtcPlugin struct {
	Opts app.WebRTCOption
}

// TODO....
type FirstConfig struct{}

func (wc *webrtcPlugin) OnEvent(event any) {
	switch event.(type) {
	case FirstConfig:
		RegisterCodecs(&wc.Opts.ME)
		i := &interceptor.Registry{}
		if len(wc.Opts.PublicIP) > 0 {
			wc.Opts.SE.SetNAT1To1IPs(wc.Opts.PublicIP, webrtc3.ICECandidateTypeHost)
		}
		if wc.Opts.PortMin > 0 && wc.Opts.PortMax > 0 {
			wc.Opts.SE.SetEphemeralUDPPortRange(wc.Opts.PortMin, wc.Opts.PortMax)
		}
		if len(wc.Opts.PublicIP) > 0 {
			wc.Opts.SE.SetNAT1To1IPs(wc.Opts.PublicIP, webrtc3.ICECandidateTypeHost)
		}
		wc.Opts.SE.SetNetworkTypes([]webrtc3.NetworkType{webrtc3.NetworkTypeUDP4, webrtc3.NetworkTypeUDP6})
		if err := webrtc3.RegisterDefaultInterceptors(&wc.Opts.ME, i); err != nil {
			panic(err)
		}
		wc.Opts.API = webrtc3.NewAPI(webrtc3.WithMediaEngine(&wc.Opts.ME),
			webrtc3.WithInterceptorRegistry(i), webrtc3.WithSettingEngine(wc.Opts.SE))
	}
}

func (wc *webrtcPlugin) Play_(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/sdp")
	streamPath := r.URL.Path[len("/webrtc/play/"):]
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var suber = WebRTCSubscriber{
		WebRTCIO: WebRTCIO{SDP: string(bytes)},
	}
	if suber.PeerConnection, err = wc.Opts.API.NewPeerConnection(webrtc3.Configuration{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	suber.OnICECandidate(func(ice *webrtc3.ICECandidate) {
		if ice != nil {
			//	suber.Info(ice.ToJSON().Candidate)
		}
	})
	if err = suber.SetRemoteDescription(webrtc3.SessionDescription{Type: webrtc3.SDPTypeOffer, SDP: suber.SDP}); err != nil {
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

func (wc *webrtcPlugin) Push_(w http.ResponseWriter, r *http.Request) {
	streamPath := r.URL.Path[len("/webrtc/push/"):]
	w.Header().Set("Content-Type", "application/sdp")
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var puber = WebRTCPublisher{WebRTCIO: WebRTCIO{SDP: string(bytes)}}
	if puber.PeerConnection, err = wc.Opts.API.NewPeerConnection(webrtc3.Configuration{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	puber.OnICECandidate(func(ice *webrtc3.ICECandidate) {
		if ice != nil {
			//puber.Info(ice.ToJSON().Candidate)
		}
	})
	if _, err = puber.AddTransceiverFromKind(webrtc3.RTPCodecTypeVideo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = puber.AddTransceiverFromKind(webrtc3.RTPCodecTypeAudio); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = webrtcObject.Publish(streamPath, &puber); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := puber.SetRemoteDescription(webrtc3.SessionDescription{Type: webrtc3.SDPTypeOffer, SDP: puber.SDP}); err != nil {
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
