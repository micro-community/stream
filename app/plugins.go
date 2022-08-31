package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/logrusorgru/aurora"
	"github.com/micro-community/stream/pubsub"
	"go-micro.dev/v4/util/log"
)

// Plugin Manager
const (
	PLUGIN_NONE       = 0      //独立插件
	PLUGIN_SUBSCRIBER = 1      //订阅者插件
	PLUGIN_PUBLISHER  = 1 << 1 //发布者插件
	PLUGIN_HOOK       = 1 << 2 //钩子插件
	PLUGIN_APP        = 1 << 3 //应用插件
)

// ListenerConfig 带有监听地址端口的插件配置类型
type ListenerConfig struct {
	ListenAddr string
}

type Plugin interface {
	// 可能的入参类型：FirstConfig 第一次初始化配置，Config 后续配置更新，SE系列（StateEvent）流状态变化事件
	OnEvent(any)
}

// PluginConfig 插件配置定义
type PluginOptions struct {
	context.Context    `json:"-"`
	context.CancelFunc `json:"-"`
	Name               string                       //插件名称
	Type               byte                         //类型
	Version            string                       //插件版本
	HotConfig          map[string]func(interface{}) //热修改配置
}

type plugin struct {
	Opts PluginOptions
	once sync.Once
}

type Option func(*PluginOptions)

// Plugins 所有的插件配置
var Plugins = make(map[string]Plugin)

// InstallPlugin 安装插件
func InstallPlugin(p PluginOptions) Plugin {

	//创建组件
	plug := &plugin{
		Opts: PluginOptions{
			Name:    "name",
			Version: "version",
		},
		once: sync.Once{},
	}

	//初始化组件
	log.Info(aurora.Green("install plugin"), aurora.BrightCyan(plug.Opts.Name), aurora.BrightBlue(plug.Opts.Version))

	return plug
}

func (p *plugin) OnEvent(any) {

}

func (p *plugin) handleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {

}

// run plugin
func (p *plugin) run() {

}

func (p *plugin) Update(opts PluginOptions) {

}

func (p *plugin) Push(streamPath string, url string, pusher pubsub.IPusher, save bool) {

}
