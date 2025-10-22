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
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
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
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
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

	// 创建资源信息
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("ShunFengParcel"),
			semconv.ServiceVersion(Version),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// APM 配置（gRPC 方式）
	apmEndpoint := "apmplus-cn-beijing.volces.com:4317" // 正确域名
	apmAppKey := "57b69ddefb13e20e8277e8d2861f4a4f"     // 替换为从火山引擎获取的真实AppKey

	var tracerProvider *trace.TracerProvider
	var meterProvider *metric.MeterProvider // 提升为函数内全局变量，确保清理函数可访问

	// 创建 OTLP Trace 导出器（gRPC 方式）
	// 创建 OTLP Trace 导出器（gRPC 方式）
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(apmEndpoint),
		otlptracegrpc.WithHeaders(map[string]string{
			"X-ByteAPM-AppKey": apmAppKey, // 修正键名：APIM →  APM
		}),
	)
	if err != nil {
		log.Fatalf("trace 导出器创建失败: %v", err)
	} else {
		log.Info("trace 导出器创建成功，已准备发送数据")
	}
	// 【Trace 导出器连接状态日志】
	if err != nil {
		log.Errorf("trace exporter 创建失败（连接服务端失败）: %v", err)
	} else {
		log.Info("trace exporter 创建成功（已连接服务端）")
		// 创建 TracerProvider
		tracerProvider = trace.NewTracerProvider(
			trace.WithBatcher(traceExporter),
			trace.WithResource(res),
		)
		otel.SetTracerProvider(tracerProvider)
	}

	// 创建 OTLP Metrics 导出器（gRPC 方式）
	metricsExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(apmEndpoint),
		otlpmetricgrpc.WithHeaders(map[string]string{
			"X-ByteAPM-AppKey": apmAppKey, // 同样修正键名
		}),
	)
	// 【Metrics 导出器连接状态日志】
	if err != nil {
		log.Errorf("metrics exporter 创建失败（连接服务端失败）: %v", err)
	} else {
		log.Info("metrics exporter 创建成功（已连接服务端）")
		// 创建 MeterProvider
		meterProvider = metric.NewMeterProvider(
			metric.WithReader(metric.NewPeriodicReader(metricsExporter)),
			metric.WithResource(res),
		)
		otel.SetMeterProvider(meterProvider)
	}

	// 返回清理函数
	return func(ctx context.Context) error {
		if tracerProvider != nil {
			if err := tracerProvider.Shutdown(ctx); err != nil {
				log.Errorf("failed to shutdown tracer provider: %v", err)
			}
		}
		if meterProvider != nil {
			if err := meterProvider.Shutdown(ctx); err != nil {
				log.Errorf("failed to shutdown meter provider: %v", err)
			}
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
