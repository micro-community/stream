package app

import (
	"context"
	"net/http"
	"time"

	webrtc3 "github.com/pion/webrtc/v3"
)

type WebRTCOption struct {
	//config.Publish
	//config.Subscribe
	Name       string
	ICEServers []string
	PublicIP   []string
	PortMin    uint16
	PortMax    uint16
	PLI        time.Duration
	ME         webrtc3.MediaEngine
	SE         webrtc3.SettingEngine
	API        *webrtc3.API
}

// PluginConfig 插件配置定义
type PluginOptions struct {
	context.Context    `json:"-"`
	context.CancelFunc `json:"-"`
	Name               string //插件名称
	Type               byte   //类型
	Version            string //插件版本
	webrtc             WebRTCOption
}

type Option func(*PluginOptions)

func WithWebRTC(opts WebRTCOption) Option {
	return func(po *PluginOptions) {
		po.webrtc = opts
	}

}

type PublishOption struct {
	PubAudio          bool
	PubVideo          bool
	KickExist         bool // 是否踢掉已经存在的发布者
	PublishTimeout    int  // 发布无数据超时
	WaitCloseTimeout  int  // 延迟自动关闭（等待重连）
	DelayCloseTimeout int  // 延迟自动关闭（无订阅时）
}

type SubscribeOption struct {
	SubAudio    bool
	SubVideo    bool
	LiveMode    bool // 实时模式：追赶发布者进度，在播放首屏后等待发布者的下一个关键帧，然后调到该帧。
	IFrameOnly  bool // 只要关键帧
	WaitTimeout int  // 等待流超时
}

type WebOption struct {
	ListenAddr    string
	ListenAddrTLS string
	CertFile      string
	KeyFile       string
	CORS          bool //是否自动添加CORS头
	UserName      string
	Password      string
	mux           *http.ServeMux
}

type PullOption struct {
	RePull          int               // 断开后自动重拉,0 表示不自动重拉，-1 表示无限重拉，高于0 的数代表最大重拉次数
	PullOnStart     bool              // 启动时拉流
	PullOnSubscribe bool              // 订阅时自动拉流
	PullList        map[string]string // 自动拉流列表，以streamPath为key，url为value
}

type ConsoleOption struct {
	Server        string //远程控制台地址
	Secret        string //远程控制台密钥
	PublicAddr    string //公网地址，提供远程控制台访问的地址，不配置的话使用自动识别的地址
	PublicAddrTLS string
}

// engine options
type EngineOption struct {
	Publish    PublishOption
	Subscribe  SubscribeOption
	HTTP       WebOption
	Console    ConsoleOption
	RTPReorder bool
	EnableAVCC bool //启用AVCC格式，rtmp协议使用
	EnableRTP  bool //启用RTP格式，rtsp、gb18181等协议使用
	LogLevel   string
}
