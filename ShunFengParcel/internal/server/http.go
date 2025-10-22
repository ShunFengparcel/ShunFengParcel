package server

import (
	pay "ShunFengParcel/api/helloworld/payment"
	v1 "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/internal/conf"
	"ShunFengParcel/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, paymentService *service.PaymentService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
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
	rount := srv.Route("/")
	rount.GET("/websocket", service.Chat)
	rount.GET("/httpwebsocket", service.HandleWebSocket)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	pay.RegisterPaymentHTTPServer(srv, paymentService)
	return srv
}
