package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func main() {
	ctx := context.Background()

	fmt.Println("========================================")
	fmt.Println("测试 OTLP 连接到 Jaeger")
	fmt.Println("========================================")

	// 创建 exporter
	fmt.Println("1. 创建 OTLP exporter...")
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint("14.103.153.242:4318"),
		otlptracehttp.WithURLPath("/v1/traces"),
	)
	if err != nil {
		log.Fatalf("❌ 创建 exporter 失败: %v", err)
	}
	fmt.Println("✅ Exporter 创建成功")

	// 创建资源
	fmt.Println("2. 创建资源信息...")
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("TestService"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("❌ 创建资源失败: %v", err)
	}
	fmt.Println("✅ 资源创建成功")

	// 创建 TracerProvider
	fmt.Println("3. 创建 TracerProvider...")
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	fmt.Println("✅ TracerProvider 创建成功")

	// 创建测试 span
	fmt.Println("4. 创建测试 span...")
	tracer := otel.Tracer("test-tracer")
	_, span := tracer.Start(ctx, "test-span")
	span.SetAttributes(
		attribute.String("test.type", "manual"),
		attribute.String("test.message", "This is a test span"),
	)
	time.Sleep(100 * time.Millisecond)
	span.End()
	fmt.Println("✅ Span 创建成功")

	// 强制刷新
	fmt.Println("5. 强制刷新数据...")
	if err := tp.ForceFlush(ctx); err != nil {
		log.Fatalf("❌ 刷新失败: %v", err)
	}
	fmt.Println("✅ 数据已刷新")

	// 关闭
	fmt.Println("6. 关闭 TracerProvider...")
	if err := tp.Shutdown(ctx); err != nil {
		log.Fatalf("❌ 关闭失败: %v", err)
	}
	fmt.Println("✅ 已关闭")

	fmt.Println("========================================")
	fmt.Println("✅ 测试完成！")
	fmt.Println("等待 10 秒后在 Jaeger UI 中查看")
	fmt.Println("Service: TestService")
	fmt.Println("Operation: test-span")
	fmt.Println("========================================")

	time.Sleep(10 * time.Second)
}
