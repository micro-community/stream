package media

import (
	"time"

	"github.com/micro-community/stream/algo"
)

type RingBuffer[T any] struct {
	*algo.Ring[T] `json:"-"`
	Size          int
	MoveCount     uint32
	LastValue     *T
}

type AVRing[T RawSlice] struct {
	RingBuffer[AVFrame[T]]
	Poll time.Duration
}
