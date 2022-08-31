package runtime

import (
	"context"

	"go-micro.dev/v4/config"
)

// Options for micro service
type Options struct {
	Config  config.Config
	Context context.Context
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
