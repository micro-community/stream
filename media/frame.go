/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:33
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-06 11:35:19
 * @FilePath: \stream\media\frame.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package media

import (
	"net"
	"time"

	"github.com/pion/rtp"
)

const (
	SUBTYPE_RAW = iota
	SUBTYPE_AVCC
	SUBTYPE_RTP
	SUBTYPE_FLV
)

type NALUSlice net.Buffers
type AudioSlice []byte

type AVCCFrame []byte   // 一帧AVCC格式的数据
type AnnexBFrame []byte // 一帧AnnexB格式数据
// RawSlice 原始切片数据
type RawSlice interface {
	~[][]byte | ~[]byte
}

type RTPFrame struct {
	rtp.Packet
}

// DecoderConfiguration for decode AV Data
type DecoderConfiguration[T RawSlice] struct {
	PayloadType byte
	AVCC        net.Buffers
	Raw         T
	FLV         net.Buffers
	Seq         int //收到第几个序列帧，用于变码率时让订阅者发送序列帧
}

// FrameBase for Media Data
type FrameBase struct {
	DeltaTime uint32    // 相对上一帧时间戳，毫秒
	AbsTime   uint32    // 绝对时间戳，毫秒
	Timestamp time.Time // 写入时间,可用于比较两个帧的先后
	Sequence  uint32    // 在一个Track中的序号
	BytesIn   int       // 输入字节数用于计算BPS
}
type AVFrame[T RawSlice] struct {
	FrameBase
	IFrame  bool
	PTS     uint32
	DTS     uint32
	AVCC    net.Buffers `json:"-"` // 打包好的AVCC格式
	RTP     []*RTPFrame `json:"-"`
	Raw     []T         `json:"-"` // 裸数据
	canRead bool
}

// Define some regular type for frame
type AudioFrame AVFrame[AudioSlice]
type VideoFrame AVFrame[NALUSlice]
type AudioDeConf DecoderConfiguration[AudioSlice]
type VideoDeConf DecoderConfiguration[NALUSlice]
type FLVFrame net.Buffers
type AudioRTP RTPFrame
type VideoRTP RTPFrame
