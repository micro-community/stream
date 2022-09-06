package webrtc

import (
	"github.com/micro-community/stream/pubsub"
	webrtcv3 "github.com/pion/webrtc/v3"
)

type WebRTCPublisher struct {
	pubsub.Publisher
	WebRTCSession
}

func (publisher *WebRTCPublisher) OnEvent(event any) {
	switch event.(type) {
	case pubsub.IPublish:
		publisher.OnICEConnectionStateChange(func(connectionState webrtcv3.ICEConnectionState) {

			switch connectionState {
			case webrtcv3.ICEConnectionStateDisconnected, webrtcv3.ICEConnectionStateFailed:
				//	publisher.Stop()
			}
		})

	}
	// publisher.Publisher.OnEvent(event)
}
