package pubsub

import "context"

type ISubscribe interface {
	IChannel
	receive(string, IChannel) error
	IsPlaying() bool
	PlayRaw()
	PlayBlock(byte)
	PlayFLV()
	Stop()
}

type TrackPlayer struct {
	context.Context    `json:"-"`
	context.CancelFunc `json:"-"`
	//Audio              PlayContext[*track.Audio, AudioSlice]
	//Video              PlayContext[*track.Video, NALUSlice]
	SkipTS     uint32 //跳过的时间戳
	FirstAbsTS uint32 //订阅起始时间戳
}

// Subscriber 订阅者实体定义
type Subscriber struct {
	Channel[SubscribeOption]
	//TrackPlayer `json:"-"`
}
