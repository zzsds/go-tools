package auth

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
	"github.com/zzsds/kratos-tools/auth"
)

var base auth.Auth

func NewAuth(a auth.Auth) {
	base = a
}

// Option is tracing option.
type Option func(*options)

type options struct {
	auth    auth.Auth
	logger  log.Logger
	header  string
	prefix  string
	exclude []Exclude
}

type Exclude map[string]string

// WithAuth with auth interface{}.
func WithAuth(auth auth.Auth) Option {
	return func(o *options) {
		o.auth = auth
	}
}

// WithLogger with recovery logger.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// WithHeader with recovery header.
func WithHeader(header string) Option {
	return func(o *options) {
		o.header = header
	}
}

// WithPrefix with recovery prefix.
func WithPrefix(prefix string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

// WithExclude with recovery exclude.
func WithExclude(exclude ...Exclude) Option {
	return func(o *options) {
		o.exclude = exclude
	}
}

// Server returns a new server middleware for OpenTelemetry.
func Server(opts ...Option) middleware.Middleware {
	options := options{
		logger: log.DefaultLogger,
		header: "Authorization",
		prefix: auth.BearerScheme,
		auth:   base,
	}
	for _, o := range opts {
		o(&options)
	}

	_ = log.NewHelper("middleware/auth", options.logger)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if info, ok := http.FromServerContext(ctx); ok {
				request := info.Request
				path := info.Request.URL.Path
				method := info.Request.Method
				// 过滤排除不验证路由
				for _, exclude := range options.exclude {
					if p, ok := exclude[method]; ok && p == path {
						return handler(ctx, req)
					}
				}
				_, ok := request.Header[options.header]
				if !ok {
					return nil, errors.Unauthorized("Unauthorized", "Header 参数名 %s 不存在", options.header)
				}
				header := request.Header.Get(options.header)
				if !strings.HasPrefix(header, options.prefix) {
					return nil, errors.Unauthorized("Unauthorized", "无效的授权头，正确应为 Bearer ")
				}
				if options.auth == nil {
					return nil, errors.DataLoss("AuthInitFail ", "认证初始化失败")
				}
				account, err := options.auth.Inspect(strings.TrimPrefix(header, options.prefix))
				if err != nil {
					return nil, errors.Unauthorized("Unauthorized", "Token 解析失败/已过期")
				}
				if request.Form == nil {
					request.Form = make(url.Values)
				}
				// 请求参数中添加 UserId
				request.Form.Add("userId", account.ID)
				// 将请求数据解析到请求结构体中
				if err := binding.BindForm(request, req); err != nil {
					return nil, errors.InvalidArgument("BindFormFail", "认证信息绑定失败")
				}
				// 将新增Form参数传递到 context.Context 中
				ctx = http.NewServerContext(ctx, info)
				// 将认证解析数据再次写入到 context.Context 中
				ctx = auth.ContextWithAccount(ctx, account)
			}
			reply, err = handler(ctx, req)
			return
		}
	}
}

// GetID 获取请求ID
func GetID(ctx context.Context) int32 {
	var id int
	account, ok := auth.AccountFromContext(ctx)
	if ok {
		id, _ = strconv.Atoi(account.ID)
	}
	return int32(id)
}

// GetType 获取请求类型【来源】
func GetType(ctx context.Context) string {
	var authType string
	account, ok := auth.AccountFromContext(ctx)
	if ok {
		authType = account.Type
	}
	return authType
}

// GetScopes 获取权限范围 【角色/权限】数组
func GetScopes(ctx context.Context) []string {
	var scopes []string
	account, ok := auth.AccountFromContext(ctx)
	if ok {
		scopes = account.Scopes
	}
	return scopes
}
