package etcd

import (
	"log"

	"github.com/go-kratos/etcd/registry"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	reg "github.com/zzsds/kratos-tools/registry"
	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type Etcd struct {
	opts *Options
	*registry.Registry
}

func Registry(opts ...Option) *Etcd {
	options := newOptions(opts...)
	cli, err := etcd.New(options.Config)
	if err != nil {
		log.Fatal(err)
	}
	return &Etcd{
		opts:     options,
		Registry: registry.New(cli, options.cliopts...),
	}
}

func (e *Etcd) ClientConn(endpoint string, opts ...transgrpc.ClientOption) grpc.ClientConnInterface {
	def := []transgrpc.ClientOption{
		transgrpc.WithEndpoint(endpoint),
		transgrpc.WithTimeout(e.opts.DialTimeout),
		transgrpc.WithDiscovery(e),
	}
	return reg.Client(def...)
}
