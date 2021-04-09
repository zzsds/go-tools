package etcd

import (
	"time"

	"github.com/go-kratos/etcd/registry"
	etcd "go.etcd.io/etcd/client/v3"
)

var DefaultAddr string = "127.0.0.1:2379"

type Options struct {
	etcd.Config
	cliopts []registry.Option
}

type Option func(*Options)

func newOptions(opts ...Option) *Options {
	options := &Options{
		Config: etcd.Config{Endpoints: []string{DefaultAddr}},
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

func WithEndpoint(endpoint ...string) Option {
	return func(o *Options) {
		o.Endpoints = endpoint
	}
}

func WithDialTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.DialTimeout = timeout
	}
}

func WithClientOption(clioption ...registry.Option) Option {
	return func(o *Options) {
		o.cliopts = clioption
	}
}
