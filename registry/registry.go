package registry

import (
	"context"
	"log"

	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

func Client(opts ...transgrpc.ClientOption) grpc.ClientConnInterface {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		opts...,
	)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
