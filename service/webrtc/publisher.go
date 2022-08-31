package webrtc

import (
	"github.com/pion/webrtc/v3"
)

type WebRTCPublisher struct {
	//Publisher
	WebRTCSession
}

type IPublisher struct{}

func (publisher *WebRTCPublisher) OnEvent(event any) {
	switch event.(type) {
	case IPublisher:
		publisher.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {

			switch connectionState {
			case webrtc.ICEConnectionStateDisconnected, webrtc.ICEConnectionStateFailed:
				//	publisher.Stop()
			}
		})

	}
	// publisher.Publisher.OnEvent(event)
}
