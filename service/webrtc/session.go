package webrtc

import (
	"github.com/pion/webrtc/v3"
)

// WebRTCSession will be changed to session
type WebRTCSession struct {
	*webrtc.PeerConnection
	SDP string
}

func (session *WebRTCSession) GetAnswer() (string, error) {
	// Sets the LocalDescription, and starts our UDP listeners
	answer, err := session.CreateAnswer(nil)
	if err != nil {
		return "", err
	}
	gatherComplete := webrtc.GatheringCompletePromise(session.PeerConnection)
	if err := session.SetLocalDescription(answer); err != nil {
		return "", err
	}
	<-gatherComplete
	return session.LocalDescription().SDP, nil
}
