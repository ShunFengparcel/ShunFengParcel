package server

import (
	v1 "ShunFengParcel/api/helloworld/v1"
	"ShunFengParcel/internal/conf"
	"ShunFengParcel/internal/service"
	"ShunFengParcel/internal/basic/config"

	_ "ShunFengParcel/internal/basic/inits"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterGreeterServer(srv, greeter)

	// 注入 Redis 客户端，避免 KuaiService.RDB 为空
	k := service.NewKuaiService(nil, config.RDB)
	v1.RegisterKuaiServer(srv, k)

	m := &service.WeiService{}
	v1.RegisterWeiServer(srv, m)
	return srv
}
