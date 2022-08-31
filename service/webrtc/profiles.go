package webrtc

import (
	"github.com/pion/webrtc/v3"
)

// init supported default RTPCodecParameters
var defaultWebrtcCodecParameters = []webrtc.RTPCodecParameters{
	{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypePCMU,
			ClockRate:    8000,
			Channels:     0,
			SDPFmtpLine:  "",
			RTCPFeedback: nil,
		},
		PayloadType: 0,
	},
	{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypePCMA,
			ClockRate:    8000,
			Channels:     0,
			SDPFmtpLine:  "",
			RTCPFeedback: nil,
		},
		PayloadType: 8,
	},
}

var supportedVideoRTCPFeedback = []webrtc.RTCPFeedback{
	{Type: "goog-remb", Parameter: ""},
	{Type: "ccm", Parameter: "fir"},
	{Type: "nack", Parameter: ""},
	{Type: "nack", Parameter: "pli"},
}

func RegisterCodecs(m *webrtc.MediaEngine) error {
	for _, codec := range defaultWebrtcCodecParameters {
		if err := m.RegisterCodec(codec, webrtc.RTPCodecTypeAudio); err != nil {
			return err
		}
	}

	for _, codec := range []webrtc.RTPCodecParameters{
		// {
		// 	RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=96", nil},
		// 	PayloadType:        97,
		// },

		// {
		// 	RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=98", nil},
		// 	PayloadType:        99,
		// },

		// {
		// 	RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=100", nil},
		// 	PayloadType:        101,
		// },

		{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:     webrtc.MimeTypeH264,
				ClockRate:    90000,
				Channels:     0,
				SDPFmtpLine:  "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42001f",
				RTCPFeedback: supportedVideoRTCPFeedback,
			},
			PayloadType: 102,
		},
		// {
		// 	RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=102", nil},
		// 	PayloadType:        121,
		// },

		{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:     webrtc.MimeTypeH264,
				ClockRate:    90000,
				Channels:     0,
				SDPFmtpLine:  "level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42001f",
				RTCPFeedback: supportedVideoRTCPFeedback,
			},
			PayloadType: 127,
		},
		// {
		// 	RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=127", nil},
		// 	PayloadType:        120,
		// },

		{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:     webrtc.MimeTypeH264,
				ClockRate:    90000,
				Channels:     0,
				SDPFmtpLine:  "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f",
				RTCPFeedback: supportedVideoRTCPFeedback,
			},
			PayloadType: 125,
		},
		// {
		// 	RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=125", nil},
		// 	PayloadType:        107,
		// },

		{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:     webrtc.MimeTypeH264,
				ClockRate:    90000,
				Channels:     0,
				SDPFmtpLine:  "level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42e01f",
				RTCPFeedback: supportedVideoRTCPFeedback,
			},
			PayloadType: 108,
		},
		// {
		// 	RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=108", nil},
		// 	PayloadType:        109,
		// },

		{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:     webrtc.MimeTypeH264,
				ClockRate:    90000,
				Channels:     0,
				SDPFmtpLine:  "level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42001f",
				RTCPFeedback: supportedVideoRTCPFeedback,
			},
			PayloadType: 127,
		},
		// {
		// 	RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=127", nil},
		// 	PayloadType:        120,
		// },

		{
			RTPCodecCapability: webrtc.RTPCodecCapability{
				MimeType:     webrtc.MimeTypeH264,
				ClockRate:    90000,
				Channels:     0,
				SDPFmtpLine:  "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=640032",
				RTCPFeedback: supportedVideoRTCPFeedback,
			},
			PayloadType: 123,
		},
		// {
		// 	RTPCodecCapability: RTPCodecCapability{"video/rtx", 90000, 0, "apt=123", nil},
		// 	PayloadType:        118,
		// },
	} {
		if err := m.RegisterCodec(codec, webrtc.RTPCodecTypeVideo); err != nil {
			return err
		}
	}
	return nil
}
