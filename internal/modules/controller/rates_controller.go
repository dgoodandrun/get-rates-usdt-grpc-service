package controller

import (
	"context"
	"get-rates-usdt-grpc-service/internal/infrastracture/metrics"
	"get-rates-usdt-grpc-service/internal/modules/service"
	pb "get-rates-usdt-grpc-service/protogen/golang/get-rates"
	"log"
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
		return nil, err
	}

	return &pb.GetRatesResponse{
		Ask:       rate.Ask,
		Bid:       rate.Bid,
		Timestamp: rate.Timestamp,
	}, nil
}

func (c *RatesController) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	status := pb.HealthCheckResponse_SERVING
	if err := c.service.HealthCheck(ctx); err != nil {
		log.Println(err)
		status = pb.HealthCheckResponse_NOT_SERVING
	}

	metrics.HealthcheckStatus.Set(float64(status))

	return &pb.HealthCheckResponse{Status: status}, nil
}
