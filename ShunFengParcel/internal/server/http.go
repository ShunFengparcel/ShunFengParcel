package server

import (
	v1 "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/internal/conf"
	"ShunFengParcel/internal/service"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, adminService *service.AdminService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(), // 添加链路追踪中间件
		),
		// 添加跨域过滤器
		http.Filter(corsFilter()),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterAdminHTTPServer(srv, adminService)
	return srv
}

// corsFilter 跨域过滤器
func corsFilter() http.FilterFunc {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			// 允许所有来源，生产环境建议指定具体域名
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// 允许的请求方法
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			// 允许的请求头
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			// 允许携带凭证
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			// 预检请求的有效期（秒）
			w.Header().Set("Access-Control-Max-Age", "86400")

			// 处理 OPTIONS 预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(stdhttp.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
