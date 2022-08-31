package webrtc

import (
	"github.com/pion/webrtc/v3"
)

type WebRTCSubscriber struct {
	//Subscriber
	WebRTCIO
	videoTrack *webrtc.TrackLocalStaticRTP
	audioTrack *webrtc.TrackLocalStaticRTP
}

func (suber *WebRTCSubscriber) OnEvent(event any) {

}
