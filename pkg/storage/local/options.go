package local

type Options struct {
	Path string
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
