/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-06 10:55:13
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 11:16:31
 * @FilePath: \stream\pubsub\stream.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package pubsub

import (
	"time"

	"github.com/micro-community/stream/algo"
	"github.com/micro-community/stream/media"
)

// Stream 流定义
type Stream struct {
	timeout           *time.Timer //当前状态的超时定时器
	actionChan        algo.SafeChan[any]
	StartTime         time.Time     //创建时间
	PublishTimeout    time.Duration //发布者无数据后超时
	DelayCloseTimeout time.Duration //发布者丢失后等待
	Path              string
	Publisher         IPublish
	State             media.StreamState
	Subscribers       []ISubscribe // 订阅者
	Tracks            media.Tracks
	AppName           string
	StreamName        string
}

func (s *Stream) Receive(event any) bool {
	//return s.actionChan.Send(event)
	return true
}
