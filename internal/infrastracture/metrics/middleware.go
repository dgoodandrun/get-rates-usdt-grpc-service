package metrics

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		res, err := handler(ctx, req)
		duration := time.Since(start).Seconds()

		status := "success"
		if err != nil {
			status = "error"
		}

		GRPCRequestsTotal.WithLabelValues(
			info.FullMethod,
			status,
		).Inc()

		GRPCRequestDuration.WithLabelValues(
			info.FullMethod,
		).Observe(duration)

		return res, err
	}
}
