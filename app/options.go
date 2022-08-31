package app

import (
	"context"
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
