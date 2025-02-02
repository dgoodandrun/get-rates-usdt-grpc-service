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

		sts := "success"
		if err != nil {
			sts = "error"
		}

		GRPCRequestsTotal.WithLabelValues(
			info.FullMethod,
			sts,
		).Inc()

		GRPCRequestDuration.WithLabelValues(
			info.FullMethod,
		).Observe(duration)

		return res, err
	}
}
