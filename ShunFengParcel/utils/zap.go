package utils

//import (
//	"ShunFengParcel/internal/server"
//	"github.com/go-kratos/kratos/contrib/log/zap/v2"
//	"github.com/go-kratos/kratos/v2"
//	"github.com/go-kratos/kratos/v2/log"
//)
//
//func Zap() {
//	logger := zap.NewLogger(zap.WithMessageKey("kratos"))
//
//	log.SetLogger(logger)
//
//	// 创建Kratos应用
//	app := kratos.New(
//		kratos.Name("your-project"),
//		kratos.Version("v1.0.0"),
//		kratos.Metadata(map[string]string{}),
//	)
//
//	// 注册HTTP服务
//	app.Handle(server.NewHTTPServer(logger))
//
//	// 启动应用
//	if err := app.Run(); err != nil {
//		log.Error(err)
//	}
//}
