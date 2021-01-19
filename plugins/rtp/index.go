package rtp

import (
	"github.com/micro-community/streaming/engine"
	"github.com/micro-community/streaming/engine/avformat"
	"github.com/micro-community/streaming/engine/util"
	"github.com/pion/rtp"
)

type RTPType int

const (
	RTP_TYPE_AUDIO RTPType = iota
	RTP_TYPE_VIDEO
	RTP_TYPE_AUDIOCONTROL
	RTP_TYPE_VIDEOCONTROL
)

type RTPPack struct {
	Type RTPType
	rtp.Packet
}

func (rt RTPType) String() string {
	switch rt {
	case RTP_TYPE_AUDIO:
		return "audio"
	case RTP_TYPE_VIDEO:
		return "video"
	case RTP_TYPE_AUDIOCONTROL:
		return "audio control"
	case RTP_TYPE_VIDEOCONTROL:
		return "video control"
	}
	return "unknow"
}

type RTP struct {
	engine.NALU
}
type RTP_PS struct {
	RTP
	rtp.Packet
	psPacket []byte
	parser   DecPSPackage
}

func (rtp *RTP_PS) PushPS(ps []byte) {
	if err := rtp.Unmarshal(ps); err != nil {
		engine.Println(err)
	}
	if len(rtp.Payload) >= 4 && util.BigEndian.Uint32(rtp.Payload) == StartCodePS {
		if rtp.psPacket != nil {
			if err := rtp.parser.Read(rtp.psPacket); err == nil {
				for _, payload := range avformat.SplitH264(rtp.parser.VideoPayload) {
					rtp.WriteNALU(rtp.Timestamp, payload)
				}
				if rtp.parser.AudioPayload != nil {
					switch rtp.parser.AudioStreamType {
					case G711A:
						rtp.AudioInfo.SoundFormat = 7
						rtp.AudioInfo.SoundRate = 8000
						rtp.AudioInfo.SoundSize = 16
						asc := rtp.AudioInfo.SoundFormat << 4
						asc = asc + 1<<1
						rtp.PushAudio(rtp.Timestamp, append([]byte{asc}, rtp.parser.AudioPayload...))
					}
				}
			} else {
				engine.Print(err)
			}
			rtp.psPacket = nil
		}
		rtp.psPacket = append(rtp.psPacket, rtp.Payload...)
	} else if rtp.psPacket != nil {
		rtp.psPacket = append(rtp.psPacket, rtp.Payload...)
	}
}
func (rtp *RTP) PushPack(pack *RTPPack) {
	switch pack.Type {
	case RTP_TYPE_AUDIO:
		payload := pack.Payload
		switch rtp.AudioInfo.SoundFormat {
		case 10:
			auHeaderLen := (int16(payload[0]) << 8) + int16(payload[1])
			auHeaderLen = auHeaderLen >> 3
			auHeaderCount := int(auHeaderLen / 2)
			var auLenArray []int
			for iIndex := 0; iIndex < int(auHeaderCount); iIndex++ {
				auHeaderInfo := (int16(payload[2+2*iIndex]) << 8) + int16(payload[2+2*iIndex+1])
				auLen := auHeaderInfo >> 3
				auLenArray = append(auLenArray, int(auLen))
			}
			startOffset := 2 + 2*auHeaderCount
			for _, auLen := range auLenArray {
				endOffset := startOffset + auLen
				addHead := []byte{0xAF, 0x01}
				rtp.PushAudio(pack.Timestamp, append(addHead, payload[startOffset:endOffset]...))
				startOffset = startOffset + auLen
			}
		case 7, 8:
			asc := rtp.AudioInfo.SoundFormat << 4
			switch {
			case rtp.AudioInfo.SoundRate >= 44000:
				asc = asc + (3 << 2)
			case rtp.AudioInfo.SoundRate >= 22000:
				asc = asc + (2 << 2)
			case rtp.AudioInfo.SoundRate >= 11000:
				asc = asc + (1 << 2)
			}
			asc = asc + 1<<1
			rtp.PushAudio(pack.Timestamp, append([]byte{asc}, payload...))
		}
	case RTP_TYPE_VIDEO:
		rtp.WriteNALU(pack.Timestamp, pack.Payload)
	}
}
