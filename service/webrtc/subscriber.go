/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:33
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-05 17:00:15
 * @FilePath: \stream\service\webrtc\subscriber.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package webrtc

import (
	"github.com/micro-community/stream/pubsub"
	webrtc3 "github.com/pion/webrtc/v3"
)

type WebRTCSubscriber struct {
	pubsub.Subscriber
	WebRTCSession
	videoTrack *webrtc3.TrackLocalStaticRTP
	audioTrack *webrtc3.TrackLocalStaticRTP
}

func (subscriber *WebRTCSubscriber) OnEvent(event any) {

}
