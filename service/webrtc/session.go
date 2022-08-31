package webrtc

import (
	"github.com/pion/webrtc/v3"
)

// WebRTCIO will be changed to session
type WebRTCIO struct {
	*webrtc.PeerConnection
	SDP string
}

func (IO *WebRTCIO) GetAnswer() (string, error) {
	// Sets the LocalDescription, and starts our UDP listeners
	answer, err := IO.CreateAnswer(nil)
	if err != nil {
		return "", err
	}
	gatherComplete := webrtc.GatheringCompletePromise(IO.PeerConnection)
	if err := IO.SetLocalDescription(answer); err != nil {
		return "", err
	}
	<-gatherComplete
	return IO.LocalDescription().SDP, nil
}
