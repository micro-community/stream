package webrtc

import (
	webrtc3 "github.com/pion/webrtc/v3"
)

type WebRTCSubscriber struct {
	//Subscriber
	WebRTCSession
	videoTrack *webrtc3.TrackLocalStaticRTP
	audioTrack *webrtc3.TrackLocalStaticRTP
}

func (subscriber *WebRTCSubscriber) OnEvent(event any) {

}
