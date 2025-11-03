package server

import (
	pay "ShunFengParcel/api/helloworld/payment"
	v1 "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/internal/conf"
	"ShunFengParcel/internal/service"

	stdhttp "net/http"

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
		// 添加CORS支持
		http.Filter(func(next stdhttp.Handler) stdhttp.Handler {
			return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
				// 设置CORS响应头
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "3600")

				// 处理预检请求
				if r.Method == "OPTIONS" {
					w.WriteHeader(stdhttp.StatusOK)
					return
				}

				next.ServeHTTP(w, r)
			})
		}),
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

	// 旧的 WebSocket 路由（保留兼容性）
	rount.GET("/websocket", service.Chat)
	rount.GET("/httpwebsocket", service.HandleWebSocket)

	// 新的通知 WebSocket 路由
	rount.GET("/ws/notifications", service.HandleNotificationWebSocket)

	v1.RegisterGreeterHTTPServer(srv, greeter)
	pay.RegisterPaymentHTTPServer(srv, paymentService)
	return srv
}
