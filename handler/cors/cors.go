package cors

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// 跨域特殊处理方式
type cors struct {
	headler http.Handler
}

func (c *cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		r.Header.Add("Access-Control-Request-Method", r.Method)
	}
	c.headler.ServeHTTP(w, r)
}

func CORS(opts ...handlers.CORSOption) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &cors{headler: handlers.CORS(opts...)(h)}
	}
}
