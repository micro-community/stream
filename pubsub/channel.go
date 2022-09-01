package pubsub

import (
	"context"
	"io"
	"net/url"
	"time"
)

// IChannel for IIO to pub/sub stream
type IChannel interface {
	IsClosed() bool
	OnEvent(any)
	Stop()
	Set(any)
	SetParentCtx(context.Context)
}

// ChannelOption
type ChannelOption interface {
	PublishOption | SubscribeOption
}

// ClientOption
type ClientOption interface {
	PullOption | PushOption
}

// Channel with Option Handled
type Channel[CO ChannelOption] struct {
	ID                 string
	Type               string
	context.Context    `json:"-"` //不要直接设置，应当通过OnEvent传入父级Context
	context.CancelFunc `json:"-"` //流关闭是关闭发布者或者订阅者
	StartTime          time.Time  //创建时间
	io.Reader          `json:"-"`
	io.Writer          `json:"-"`
	io.Closer          `json:"-"`
	Args               url.Values
	ChannelOption      *CO      `json:"-"`
	Specs              IChannel `json:"-"`
}
