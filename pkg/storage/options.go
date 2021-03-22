package storage

import (
	"context"
)

// Options contains configuration for the Storage
type Options struct {
	// Context should contain all implementation specific options, using context.WithValue.
	Context context.Context
}

// Option sets values in Options
type Option func(o *Options)

// WithContext sets the Storages context, for any extra configuration
func WithContext(c context.Context) Option {
	return func(o *Options) {
		o.Context = c
	}
}
