package etcd

import (
	"log"

	"github.com/go-kratos/etcd/registry"
	etcd "go.etcd.io/etcd/client/v3"
)

var DefaultAddr string = "127.0.0.1:2379"

func Registry(endpoints ...string) *registry.Registry {
	if len(endpoints) < 1 {
		endpoints = []string{DefaultAddr}
	}
	cli, err := etcd.New(etcd.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	return registry.New(cli)
}
