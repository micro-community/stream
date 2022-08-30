package pubsub

import (
	"context"
)

type IIO interface {
	IsClosed() bool
	OnEvent(any)
	Stop()
	SetIO(any)
	SetParentCtx(context.Context)
}

type ISubscriber interface {
	IIO
	receive(string, IIO) error
	GetIO()
	GetConfig()
	IsPlaying() bool
	PlayRaw()
	PlayBlock(byte)
	PlayFLV()
	Stop()
}
