package main

import (
	"ShunFengParcel/internal/conf"
	"context"
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	_ "github.com/go-kratos/kratos/v2/encoding/json"
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
	// 默认配置路径，支持多种启动方式
	defaultConf := "../../configs"
	// 如果从项目根目录启动（使用编译后的二进制），使用 ./configs
	if _, err := os.Stat("./configs"); err == nil {
		defaultConf = "./configs"
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

func initTracerProvider() func(context.Context) error {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("14.103.153.242:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	// 创建资源信息
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("ShunFengParcel"),
			semconv.ServiceVersion("1.0.0"),
			semconv.ServiceInstanceID(id),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}
	// 创建 TracerProvider
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter), // 使用批处理导出
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()), // 总是采样
	)
	otel.SetTracerProvider(tracerProvider)

	log.Infof("✅ OTLP Tracer initialized")
	log.Infof("   Endpoint: 14.103.153.242:4317 (OpenTelemetry Collector)")
	log.Infof("   Service Name: ShunFengParcel")
	log.Infof("   Collector will forward traces to Jaeger")

	// 返回清理函数
	return func(ctx context.Context) error {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Errorf("failed to shutdown tracer provider: %v", err)
			return err
		}
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

	// 初始化 OpenTelemetry
	shutdown := initTracerProvider()
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			// 忽略关闭错误
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

	// 启动服务（阻塞执行）
	if err := app.Run(); err != nil {
		panic(err)
	}
}
