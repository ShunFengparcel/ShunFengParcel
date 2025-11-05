package server

import (
	v1 "ShunFengParcel/api/helloworld/v1"

	"ShunFengParcel/internal/conf"
	"ShunFengParcel/internal/service"

	_ "ShunFengParcel/internal/basic/inits"

	"github.com/go-kratos-ecosystem/components/v2/middleware/cors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, kuai *service.KuaiService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			kuai.CourierAuthMiddleware(),
			cors.Cors(
				cors.AllowedOrigins("*"), // 允许所有来源，生产环境建议指定具体域名
				cors.AllowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS"),
				cors.AllowedHeaders("Content-Type", "Authorization", "X-Courier-ID"),
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
    v1.RegisterKuaiHTTPServer(srv, kuai)
    m := &service.WeiService{}
    v1.RegisterWeiHTTPServer(srv, m)

    // 注册 WebSocket 实时位置订阅：/ws/order/{orderId}
    hub := NewWSHub()
    hub.RegisterRoutes(srv)

    return srv
}
