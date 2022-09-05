/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-05 17:14:18
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-05 17:41:10
 * @FilePath: \stream\track\define.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package track

import (
	"github.com/micro-community/stream/media"
	"github.com/micro-community/stream/util"
)

// Media 基础媒体Track类
type Media[T media.RawSlice] struct {
	media.TrackBase
	media.AVRing[T]
	SampleRate           uint32
	DecoderConfiguration media.DecoderConfiguration[T] `json:"-"` //H264(SPS、PPS) H265(VPS、SPS、PPS) AAC(config)
	// util.BytesPool                               //无锁内存池，用于发布者（在同一个协程中）复用小块的内存，通常是解包时需要临时使用
	rtpSequence uint16            //用于生成下一个rtp包的序号
	orderQueue  []*media.RTPFrame //rtp包的缓存队列，用于乱序重排
	lastSeq     uint16            //上一个收到的序号，用于乱序重排
	lastSeq2    uint16            //记录上上一个收到的序列号
	//流速控制
}

func (av *Media[T]) ReadRing() *media.AVRing[T] {
	return util.Clone(av.AVRing)
}
