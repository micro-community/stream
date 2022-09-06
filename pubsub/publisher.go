/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:33
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 12:20:50
 * @FilePath: \stream\pubsub\publisher.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package pubsub

import "github.com/micro-community/stream/media"

// Publish of Publisher
type IPublish interface {
	IChannel[PublishOption]
	//get option
	// implemented in IChannel
	//GetOption() *PublishOption
	//receive stream from channe
	//Get the channel instance for publishing
	// implemented in IChannel
	//GetChannel() *Channel[PublishOption]

	//get track
	getAudioTrack() media.AudioTrack
	getVideoTrack() media.VideoTrack

	WriteAVCCVideo(ts uint32, frame media.AVCCFrame)
	WriteAVCCAudio(ts uint32, frame media.AVCCFrame)
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

func (p *Publisher) getAudioTrack() media.AudioTrack {
	return p.AudioTrack
}
func (p *Publisher) getVideoTrack() media.VideoTrack {
	return p.VideoTrack
}

//WriteAVCCVideo Data Frame
func (p *Publisher) WriteAVCCVideo(ts uint32, frame media.AVCCFrame) {

}

func (p *Publisher) WriteAVCCAudio(ts uint32, frame media.AVCCFrame) {

}
