package main

import (
	"context"
	"flag"
	"os"
	"time"

	"ShunFengParcel/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	// 默认配置路径，支持多种运行方式
	defaultConf := "../../configs"
	if _, err := os.Stat(defaultConf); os.IsNotExist(err) {
		// 如果 ../../configs 不存在，尝试 ./configs
		if _, err := os.Stat("./configs"); err == nil {
			defaultConf = "./configs"
		}
	}
	flag.StringVar(&flagconf, "conf", defaultConf, "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

// initJaeger 初始化 Jaeger 链路追踪（使用 OTLP HTTP 协议）
func initJaeger() func(context.Context) error {
	ctx := context.Background()

	// Jaeger OTLP HTTP 端点（新端口：8318）
	jaegerEndpoint := "14.103.153.242:8318"

	log.Infof("========================================")
	log.Infof("正在初始化 Jaeger 链路追踪...")
	log.Infof("Jaeger 服务器: %s", jaegerEndpoint)
	log.Infof("协议: OTLP HTTP")
	log.Infof("服务名称: ShunFengParcel")
	log.Infof("========================================")

	// 创建 OTLP HTTP exporter（Jaeger v2 支持）
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(jaegerEndpoint),
		otlptracehttp.WithURLPath("/v1/traces"),
	)
	if err != nil {
		log.Errorf("❌ Jaeger 连接失败: %v", err)
		log.Errorf("   请检查:")
		log.Errorf("   1. Jaeger 服务器是否运行")
		log.Errorf("   2. 网络连接是否正常")
		log.Errorf("   3. 端口 4318 是否开放")
		log.Errorf("   服务将继续运行，但不会上报链路数据")
		return func(ctx context.Context) error { return nil }
	}

	log.Infof("✅ OTLP Exporter 创建成功")

	// 创建资源信息
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("ShunFengParcel"),
			semconv.ServiceVersion(Version),
		),
	)
	if err != nil {
		log.Errorf("⚠️  资源创建失败: %v", err)
	} else {
		log.Infof("✅ 资源信息创建成功")
	}

	// 创建 TracerProvider
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	log.Infof("✅ TracerProvider 已设置")

	// 发送测试 span 验证连接
	log.Infof("正在发送测试 span 验证连接...")
	testTracer := otel.Tracer("ShunFengParcel")
	testCtx, testSpan := testTracer.Start(ctx, "jaeger-connection-test")
	testSpan.SetAttributes(
		attribute.String("test.type", "connection"),
		attribute.String("test.status", "success"),
		attribute.String("message", "Jaeger connection test successful"),
	)
	time.Sleep(100 * time.Millisecond) // 模拟一些工作
	testSpan.End()

	// 强制刷新，确保测试 span 被发送
	if err := tp.ForceFlush(testCtx); err != nil {
		log.Errorf("⚠️  测试 span 刷新失败: %v", err)
	} else {
		log.Infof("✅ 测试 span 已发送")
	}

	log.Infof("========================================")
	log.Infof("🎉 Jaeger 链路追踪初始化成功！")
	log.Infof("   - 采样率: 100%% (所有请求都会被追踪)")
	log.Infof("   - 查看链路: http://14.103.153.242:16686")
	log.Infof("   - 服务名称: ShunFengParcel")
	log.Infof("   - 测试 span: jaeger-connection-test")
	log.Infof("========================================")
	log.Infof("💡 提示: 等待 5-10 秒后在 Jaeger UI 中查看测试链路")
	log.Infof("========================================")

	// 返回清理函数
	return func(ctx context.Context) error {
		log.Infof("正在关闭 Jaeger TracerProvider...")
		if err := tp.Shutdown(ctx); err != nil {
			log.Errorf("TracerProvider 关闭失败: %v", err)
			return err
		}
		log.Infof("✅ Jaeger TracerProvider 已关闭")
		return nil
	}
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	// 初始化 Jaeger（Jaeger v2 已启动，使用 OTLP 协议）
	shutdown := initJaeger()
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Errorf("failed to shutdown tracer: %v", err)
		}
	}()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
