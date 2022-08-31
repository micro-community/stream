package app

import "context"

// PluginConfig 插件配置定义
type PluginOptions struct {
	context.Context    `json:"-"`
	context.CancelFunc `json:"-"`
	Name               string //插件名称
	Type               byte   //类型
	Version            string //插件版本
}

type Option func(*PluginOptions)
