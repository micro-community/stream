package app

import (
	"time"
)

// Parameters for engine
type Parameters struct {
	EnableWaitStream bool
	RingSize         int
	EnableAudio      bool
	EnableVideo      bool
	PublishTimeout   time.Duration
}

// ExtendInfo for extension
type ExtendInfo struct {
	Version   *string
	StartTime time.Time //启动时间
	Params    *Parameters
}

// Settings for engine
var (
	Config = &Parameters{true, 10, false, true, time.Minute}

	// ConfigRaw 配置信息的原始数据
	ConfigRaw []byte
	// Version 引擎版本号
	Version string
	// AppInfo 引擎信息
	AppInfo = &ExtendInfo{&Version, time.Now(), Config}
)

type Configuration struct {
	Params Parameters
	global struct {
		console struct {
			secrets string
		}
	}

	webrtc struct {
		enable bool
	}

	rtsp struct {
		enable bool
		pull   struct {
			pullstart bool
			pulllist  map[string]string
		}
	}
}
