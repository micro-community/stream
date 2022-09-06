/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:33
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 12:22:34
 * @FilePath: \stream\pubsub\channel.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package pubsub

import (
	"context"
	"io"
	"net/url"
	"time"
)

// Client for Pubsub Options
type Client[ClientOpt ClientOption] struct {
	Option         *ClientOpt //客户端连接选项
	StreamPath     string     // 本地流标识
	RemoteURL      string     // 远程服务器地址（用于推拉）
	ReConnectCount int        //重连次数
}

// IChannel for IIO to pub/sub stream
type IChannel[ChannelOpt ChannelOption] interface {
	IsClosed() bool
	OnEvent(any)
	Stop()
	Set(any)
	SetParentCtx(context.Context)
	//receive
	receive(streamPathUrl string, channel IChannel[ChannelOpt], opt *ChannelOpt) error
}

// Channel with Option Handled
type Channel[ChannelOpt ChannelOption] struct {
	ID                 string
	Type               string
	Args               url.Values
	StartTime          time.Time            //创建时间
	Stream             *Stream              `json:"-"`
	ChannelOption      *ChannelOpt          `json:"-"`
	Channel            IChannel[ChannelOpt] `json:"-"`
	context.Context    `json:"-"`           //不要直接设置，应当通过OnEvent传入父级Context
	context.CancelFunc `json:"-"`           //流关闭是关闭发布者或者订阅者
	io.Reader          `json:"-"`
	io.Writer          `json:"-"`
	io.Closer          `json:"-"`
}

func (ch *Channel[ChannelOpt]) IsClosed() bool {
	return ch.Err() != nil
}

// handle Event from Stream
func (ch *Channel[ChannelOpt]) OnEvent(any) {

}

func (ch *Channel[ChannelOpt]) Stop() {

}

// Set Writer、Reader、Closer
func (ch *Channel[ChannelOpt]) Set(operator any) {
	if v, ok := operator.(io.Closer); ok {
		ch.Closer = v
		return
	}
	if v, ok := operator.(io.Reader); ok {
		ch.Reader = v
		return
	}
	if v, ok := operator.(io.Writer); ok {
		ch.Writer = v
		return
	}
}

func (ch *Channel[ChannelOpt]) SetParentCtx(context.Context) {

}

func (ch *Channel[ChannelOpt]) GetChannel() *Channel[ChannelOpt] {
	return ch
}

// receive from channel
func (ch *Channel[ChannelOpt]) receive(streamPathUrl string, channel IChannel[ChannelOpt], opt *ChannelOpt) error {
	return nil
}
