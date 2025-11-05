package server

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// CourierIDKey is the key for courier id in context.
type CourierIDKey struct{}

// CourierAuthMiddleware is a middleware to get courier id from header.
func CourierAuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				courierIDStr := header.RequestHeader().Get("X-Courier-ID")
				if courierID, err := strconv.ParseInt(courierIDStr, 10, 64); err == nil {
					ctx = context.WithValue(ctx, CourierIDKey{}, courierID)
				}
			}
			return handler(ctx, req)
		}
	}
}
