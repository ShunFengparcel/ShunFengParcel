package server

import (
	v1 "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/internal/conf"
	"ShunFengParcel/internal/service"
	"ShunFengParcel/internal/basic/config"

	_ "ShunFengParcel/internal/basic/inits"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos-ecosystem/components/v2/middleware/cors"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			cors.Cors(
				cors.AllowedOrigins("*"), // 允许所有来源，生产环境建议指定具体域名
				cors.AllowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"),
				cors.AllowedHeaders("*"), // 允许所有请求头
			),
		),
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
	// 修复 KuaiService 未注入 Redis 客户端导致的空指针
	k := service.NewKuaiService(nil, config.RDB)
	v1.RegisterKuaiHTTPServer(srv, k)
	m := &service.WeiService{}
	v1.RegisterWeiHTTPServer(srv, m)

	return srv
}
