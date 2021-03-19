package casbin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/zzsds/kratos-tools/auth"
)

var base auth.Auth

func NewAuth(a auth.Auth) {
	base = a
}

// Option is tracing option.
type Option func(*options)

type options struct {
	logger log.Logger
}

// WithLogger with recovery logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// Server returns a new server middleware for OpenTelemetry.
func Server(opts ...Option) middleware.Middleware {
	options := options{
		logger: log.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}

	_ = log.NewHelper("middleware/auth", options.logger)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				path      string
				method    string
				params    string
				component string
			)
			if info, ok := http.FromServerContext(ctx); ok {
				component = "HTTP"
				// request := info.Request
				path = info.Request.RequestURI
				method = info.Request.Method
				params = info.Request.Form.Encode()

				_, _ = auth.AccountFromContext(ctx)

			} else if info, ok := grpc.FromServerContext(ctx); ok {
				component = "gRPC"
				path = info.FullMethod
				method = "POST"
			}
			_, _, _, _ = path, method, params, component
			fmt.Println()
			reply, err = handler(ctx, req)
			return
		}
	}
}

func GetID(ctx context.Context) int32 {
	var id int
	account, ok := auth.AccountFromContext(ctx)
	if ok {
		id, _ = strconv.Atoi(account.ID)
	}
	return int32(id)
}
