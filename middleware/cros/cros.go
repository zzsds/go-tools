package cors

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
)

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

	_ = log.NewHelper("middleware/cros", options.logger)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if info, ok := transhttp.FromServerContext(ctx); ok {
				request := info.Request
				origin := request.Header.Get("Origin")
				info.Response.Header().Set("Access-Control-Allow-Origin", origin)
				info.Response.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Authorization")
				info.Response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
				info.Response.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
				info.Response.Header().Set("Access-Control-Allow-Credentials", "true")

				ctx = transhttp.NewServerContext(ctx, info)
				// 放行所有OPTIONS方法
				if request.Method == http.MethodOptions {
					info.Response.WriteHeader(http.StatusNoContent)
					// OPTIONS 直接返回
					return nil, errors.Error(0, "", "")
				}
			}
			reply, err = handler(ctx, req)
			return
		}
	}
}
