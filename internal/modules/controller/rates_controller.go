package controller

import (
	"context"
	"get-rates-usdt-grpc-service/internal/infrastructure/metrics"
	"get-rates-usdt-grpc-service/internal/modules/service"
	pb "get-rates-usdt-grpc-service/protogen/golang/get-rates"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RatesController struct {
	pb.UnimplementedRatesServiceServer
	service service.RatesService
}

func NewRatesController(s service.RatesService) *RatesController {
	return &RatesController{service: s}
}

func (c *RatesController) GetRates(ctx context.Context, req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	rate, err := c.service.GetCurrentRate(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get rates: %v", err)
	}

	return &pb.GetRatesResponse{
		Ask:       rate.Ask,
		Bid:       rate.Bid,
		Timestamp: rate.Timestamp,
	}, nil
}

func (c *RatesController) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	sts := pb.HealthCheckResponse_SERVING

	if err := c.service.HealthCheck(ctx); err != nil {
		sts = pb.HealthCheckResponse_NOT_SERVING
		return &pb.HealthCheckResponse{
			Status: sts,
		}, status.Errorf(codes.Unavailable, "Service is unhealthy")
	}

	metrics.HealthcheckStatus.Set(float64(sts))

	return &pb.HealthCheckResponse{
		Status: sts,
	}, nil
}
