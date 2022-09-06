/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:33
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 09:27:38
 * @FilePath: \stream\service\webrtc\session.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package webrtc

import (
	webrtc3 "github.com/pion/webrtc/v3"
)

// WebRTCSession will be changed to session
type WebRTCSession struct {
	*webrtc3.PeerConnection
	SDP string
}

func (session *WebRTCSession) GetAnswer() (string, error) {
	// Sets the LocalDescription, and starts our UDP listeners
	answer, err := session.CreateAnswer(nil)
	if err != nil {
		return "", err
	}
	gatherComplete := webrtc3.GatheringCompletePromise(session.PeerConnection)
	if err := session.SetLocalDescription(answer); err != nil {
		return "", err
	}
	<-gatherComplete
	return session.LocalDescription().SDP, nil
}
