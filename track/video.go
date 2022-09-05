/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-05 17:08:58
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-05 17:37:12
 * @FilePath: \stream\track\video.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package track

import (
	"github.com/micro-community/stream/algo"
	"github.com/micro-community/stream/codecs"
	"github.com/micro-community/stream/media"
)

// Video for video track infos
type Video struct {
	Media[media.NALUSlice]
	CodecID     codecs.VideoCodecID
	IDRing      *algo.Ring[media.AVFrame[media.NALUSlice]] `json:"-"` //最近的关键帧位置，首屏渲染
	SPSInfo     codecs.SPSInfo
	GOP         int  //关键帧间隔
	naluLenSize int  //avcc格式中表示nalu长度的字节数，通常为4
	idrCount    int  //缓存中包含的idr数量
	dcChanged   bool //解码器配置是否改变了，一般由于变码率导致
	dtsEst      *media.DTSEstimator
}

func (v *Video) GetDecConfSeq() int {
	return v.DecoderConfiguration.Seq
}

func (vt *Video) ReadRing() *media.AVRing[media.NALUSlice] {
	vr := vt.Media.ReadRing()
	vr.Ring = vt.IDRing
	return vr
}
