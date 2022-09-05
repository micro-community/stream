package track

import (
	"github.com/micro-community/stream/codecs"
	"github.com/micro-community/stream/media"
)

var adcflv1 = []byte{codecs.FLV_TAG_TYPE_AUDIO, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0}
var adcflv2 = []byte{0, 0, 0, 15}

// Audio for track
type Audio struct {
	Media[media.AudioSlice]
	CodecID    codecs.AudioCodecID
	Channels   byte
	SampleSize byte
	AVCCHead   []byte // 音频包在AVCC格式中，AAC会有两个字节，其他的只有一个字节
	// Profile:
	// 0: Main profile
	// 1: Low Complexity profile(LC)
	// 2: Scalable Sampling Rate profile(SSR)
	// 3: Reserved
	Profile byte
}

func (a *Audio) GetDecConfSeq() int {
	return a.DecoderConfiguration.Seq
}
