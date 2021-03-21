package yunpian

// Options ...
type Options struct {
	Key   string
	Debug bool
}

type Option func(*Options)

func WithKey(key string) Option {
	return func(o *Options) {
		o.Key = key
	}
}

func WithDebug(debug bool) Option {
	return func(o *Options) {
		o.Debug = debug
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Debug: false,
	}

	for _, o := range opts {
		o(&opt)
	}
	return opt
}
