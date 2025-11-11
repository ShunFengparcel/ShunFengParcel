package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Server 日志中间件（服务端）
func Server(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				kind      string
				operation string
			)
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}

			startTime := time.Now()
			reply, err := handler(ctx, req)
			duration := time.Since(startTime)

			level := log.LevelInfo
			if err != nil {
				level = log.LevelError
			}

			_ = log.WithContext(ctx, logger).Log(level,
				"kind", "server",
				"component", kind,
				"operation", operation,
				"args", extractArgs(req),
				"duration", duration.Milliseconds(),
				"error", err,
			)

			return reply, err
		}
	}
}

// Client 日志中间件（客户端）
func Client(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				kind      string
				operation string
			)
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}

			startTime := time.Now()
			reply, err := handler(ctx, req)
			duration := time.Since(startTime)

			level := log.LevelInfo
			if err != nil {
				level = log.LevelError
			}

			_ = log.WithContext(ctx, logger).Log(level,
				"kind", "client",
				"component", kind,
				"operation", operation,
				"args", extractArgs(req),
				"duration", duration.Milliseconds(),
				"error", err,
			)

			return reply, err
		}
	}
}

// extractArgs 提取请求参数（简化版，避免敏感信息）
func extractArgs(req interface{}) string {
	if req == nil {
		return ""
	}
	return fmt.Sprintf("%T", req)
}
