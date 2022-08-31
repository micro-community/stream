package webrtc

import (
	"github.com/pion/webrtc/v3"
)

type WebRTCPublisher struct {
	//Publisher
	WebRTCIO
}

type IPublisher struct{}

func (puber *WebRTCPublisher) OnEvent(event any) {
	switch event.(type) {
	case IPublisher:
		puber.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {

			switch connectionState {
			case webrtc.ICEConnectionStateDisconnected, webrtc.ICEConnectionStateFailed:
				//	puber.Stop()
			}
		})

	}
	// puber.Publisher.OnEvent(event)
}
