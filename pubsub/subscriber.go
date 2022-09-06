/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:33
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 09:51:41
 * @FilePath: \stream\pubsub\subscriber.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package pubsub

import "github.com/micro-community/stream/media"

type ISubscribe interface {
	IChannel
	IsPlaying() bool
	PlayRaw()
	PlayBlock(byte)
	PlayFLV()
	Stop()
}

// Subscriber 订阅者实体定义
type Subscriber struct {
	Channel[SubscribeOption]
	TrackPlayer `json:"-"`
}

func (s *Subscriber) OnEvent(event any) {

}

func (s *Subscriber) AddTrack(t media.Track) bool {

	return false
}

func (s *Subscriber) IsPlaying() bool {
	return s.TrackPlayer.Context != nil && s.TrackPlayer.Err() == nil
}

func (s *Subscriber) PlayRaw() {

}

func (s *Subscriber) PlayFLV() {

}

func (s *Subscriber) PlayRTP() {

}

//PlayBlock 阻塞式读取数据
func (s *Subscriber) PlayBlock(subType byte) {

}
