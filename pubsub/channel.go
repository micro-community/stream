package pubsub

import "context"

type IChannel interface {
	IsClosed() bool
	OnEvent(any)
	Stop()
	Set(any)
	SetParentCtx(context.Context)
}
