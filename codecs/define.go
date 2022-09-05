/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-05 17:11:56
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-05 17:12:36
 * @FilePath: \stream\codecs\define.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package codecs

type AudioCodecID byte
type VideoCodecID byte

const (
	ADTS_HEADER_SIZE              = 7
	CodecID_AAC      AudioCodecID = 0xA
	CodecID_PCMA     AudioCodecID = 7
	CodecID_PCMU     AudioCodecID = 8
	CodecID_H264     VideoCodecID = 7
	CodecID_H265     VideoCodecID = 0xC
)
