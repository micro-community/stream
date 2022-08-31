package app

import (
	"net/http"
	"sync"

	"github.com/logrusorgru/aurora"
	"github.com/micro-community/stream/pubsub"
	"go-micro.dev/v4/util/log"
)

type Plugin interface {
	// 可能的入参类型：FirstConfig 第一次初始化配置，Config 后续配置更新，SE系列（StateEvent）流状态变化事件
	OnEvent(any)
	Publish(streamPath string, pub pubsub.IPublisher) error
	Subscribe(streamPath string, sub pubsub.ISubscriber) error
}

type plugin struct {
	Opts PluginOptions
	once sync.Once
}

// Plugins 所有的插件配置
var Plugins = make(map[string]Plugin)

// Install 安装功能组件
func Install(opts ...Option) Plugin {

	pluginOptions := PluginOptions{
		Name:    "DefaultName",
		Version: "0.0.1",
	}
	for _, o := range opts {
		o(&pluginOptions)
	}
	//创建组件
	plug := &plugin{
		Opts: pluginOptions,
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

func (p *plugin) Publish(streamPath string, pub pubsub.IPublisher) error {

	return nil
}

func (p *plugin) Subscribe(streamPath string, sub pubsub.ISubscriber) error {

	return nil
}
