/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:33
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 11:17:00
 * @FilePath: \stream\pubsub\publisher.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package pubsub

import "github.com/micro-community/stream/media"

// Publish of Publisher
type IPublish interface {
	IChannel
	//get option
	GetOption() *PublishOption
	//receive stream from channel
	receive(string, IChannel, PublishOption) error
	//Get the channel instance for publishing
	GetChannel() *Channel[PublishOption]

	//get track
	getAudioTrack() media.AudioTrack
	getVideoTrack() media.VideoTrack
}

// Publisher 发布者定义
type Publisher struct {
	Channel[PublishOption]
	media.AudioTrack `json:"-"`
	media.VideoTrack `json:"-"`
}

func (p *Publisher) GetOption() *PublishOption {

	return nil
}

func (p *Publisher) Stop() {
	p.Channel.Stop()
	p.Stream.Receive(media.ACTION_PUBLISHLOST)
}
