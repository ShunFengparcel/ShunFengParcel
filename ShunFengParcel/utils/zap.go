package utils

//import (
//	"ShunFengParcel/internal/server"
//	"github.com/go-kratos/kratos/v2"
//	"github.com/go-kratos/kratos/v2/log"
//	"zap"
//)
//
//func Zap() {
//	logger, err := zap.NewLogger(zap.WithName("kratos"))
//	if err != nil {
//		panic(err)
//	}
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
