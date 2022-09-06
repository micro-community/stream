/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:33
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 09:57:30
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
type IChannel interface {
	IsClosed() bool
	OnEvent(any)
	Stop()
	Set(any)
	SetParentCtx(context.Context)
}

// Channel with Option Handled
type Channel[ChannelOpt ChannelOption] struct {
	ID                 string
	Type               string
	context.Context    `json:"-"` //不要直接设置，应当通过OnEvent传入父级Context
	context.CancelFunc `json:"-"` //流关闭是关闭发布者或者订阅者
	StartTime          time.Time  //创建时间
	io.Reader          `json:"-"`
	io.Writer          `json:"-"`
	io.Closer          `json:"-"`
	Args               url.Values
	ChannelOption      *ChannelOpt `json:"-"`
	Specs              IChannel    `json:"-"`
}

func (ch *Channel[ChannelOpt]) IsClosed() bool {
	return ch.Err() != nil
}

// handle Event from Stream
func (ch *Channel[ChannelOpt]) OnEvent(any) {

}

func (ch *Channel[ChannelOpt]) Stop() {

}

func (ch *Channel[ChannelOpt]) Set(any) {

}

func (ch *Channel[ChannelOpt]) SetParentCtx(context.Context) {

}

// receive from channel
func (ch *Channel[ChannelOpt]) receive(streamPath string, channel IChannel, sOpt *SubscribeOption) {

}
