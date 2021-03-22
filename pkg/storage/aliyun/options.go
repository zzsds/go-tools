package aliyun

type Options struct {
	Path      string
	Bucket    string
	AclType   string
	Endpoint  string
	AccessKey string
	SecretKey string
	ClassType string
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func WithPath(path string) Option {
	return func(o *Options) {
		o.Path = path
	}
}

func WithBucket(bucket string) Option {
	return func(o *Options) {
		o.Bucket = bucket
	}
}

func WithEndpoint(endpoint string) Option {
	return func(o *Options) {
		o.Endpoint = endpoint
	}
}

func WithClassType(class string) Option {
	return func(o *Options) {
		o.ClassType = class
	}
}

func WithAccessKey(key string) Option {
	return func(o *Options) {
		o.AccessKey = key
	}
}

func WithSecretKey(key string) Option {
	return func(o *Options) {
		o.SecretKey = key
	}
}
