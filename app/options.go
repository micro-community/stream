/*
 * @Author: Edward crazybber@outlook.com
 * @Date: 2022-09-02 12:47:32
 * @LastEditors: Edward crazybber@outlook.com
 * @LastEditTime: 2022-09-05 16:53:11
 * @FilePath: \stream\app\options.go
 * @Description: code content
 * Copyright (c) 2022 by Edward crazybber@outlook.com, All Rights Reserved.
 */
package app

import (
	"context"
	"net/http"
	"time"

	"github.com/micro-community/stream/pubsub"
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

type ConsoleOption struct {
	Server        string //远程控制台地址
	Secret        string //远程控制台密钥
	PublicAddr    string //公网地址，提供远程控制台访问的地址，不配置的话使用自动识别的地址
	PublicAddrTLS string
}

// engine options
type EngineOption struct {
	Publish    pubsub.PublishOption
	Subscribe  pubsub.SubscribeOption
	HTTP       WebOption
	Console    ConsoleOption
	RTPReorder bool
	EnableAVCC bool //启用AVCC格式，rtmp协议使用
	EnableRTP  bool //启用RTP格式，rtsp、gb18181等协议使用
	LogLevel   string
}
