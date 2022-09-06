/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:33
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 11:02:49
 * @FilePath: \stream\media\track.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package media

import (
	"time"

	"github.com/micro-community/stream/algo"
	"github.com/pion/rtp"
)

type TimelineData[T any] struct {
	Timestamp time.Time
	Value     T
}

// TrackBase 基础Track类
type TrackBase struct {
	Name    string
	Stream  IStream `json:"-"` //所属Stream
	ts      time.Time
	bytes   int
	frames  int
	BPS     int
	FPS     int
	RawPart []int               // 裸数据片段用于UI上显示
	RawSize int                 // 裸数据长度
	BPSs    []TimelineData[int] // 10s码率统计
	FPSs    []TimelineData[int] // 10s帧率统计
}

type Track interface {
	GetBase() *TrackBase
	LastWriteTime() time.Time
	SnapForJson()
}

// AVTrack for audio and video
type AVTrack interface {
	Track
	Attach()
	Detach()
	WriteAVCC(ts uint32, frame AVCCFrame) //写入AVCC格式的数据
	WriteRTP([]byte)
	WriteRTPPack(*rtp.Packet)
	Flush()
}

// VideoTrack for Video
type VideoTrack interface {
	AVTrack
	GetDecoderConfiguration() DecoderConfiguration[NALUSlice]
	CurrentFrame() *AVFrame[NALUSlice]
	PreFrame() *AVFrame[NALUSlice]
	WriteSlice(NALUSlice)
	WriteAnnexB(uint32, uint32, AnnexBFrame)
}

type AudioTrack interface {
	AVTrack
	GetDecoderConfiguration() DecoderConfiguration[AudioSlice]
	CurrentFrame() *AVFrame[AudioSlice]
	PreFrame() *AVFrame[AudioSlice]
	WriteSlice(AudioSlice)
	WriteADTS([]byte)
}

// Tracks for streams
type Tracks struct {
	algo.Map[string, Track]
}
