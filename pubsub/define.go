package pubsub

import (
	"context"

	"github.com/micro-community/stream/media"
	"github.com/micro-community/stream/track"
)

type ReadType[T media.RawSlice] interface {
	GetDecConfSeq() int
	ReadRing() *media.AVRing[T]
}

type PlayContext[T ReadType[R], R media.RawSlice] struct {
	Track   T
	ring    *media.AVRing[R]
	confSeq int
	Frame   **media.AVFrame[R]
}
type TrackPlayer struct {
	context.Context    `json:"-"`
	context.CancelFunc `json:"-"`
	Audio              PlayContext[*track.Audio, media.AudioSlice]
	Video              PlayContext[*track.Video, media.NALUSlice]
	SkipTS             uint32 //跳过的时间戳
	FirstAbsTS         uint32 //订阅起始时间戳
}
