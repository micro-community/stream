package engine

import (
	"context"

	"github.com/micro/go-micro/v2/config"
)

// Options for micro service
type Options struct {
	Config  config.Config
	Context context.Context
}

func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
