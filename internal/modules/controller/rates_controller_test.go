package controller

import (
	"context"
	"get-rates-usdt-grpc-service/internal/models"
	"get-rates-usdt-grpc-service/internal/modules/service/mocks"
	pb "get-rates-usdt-grpc-service/protogen/golang/get-rates"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRatesController_GetRates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockRatesService(ctrl)
	controller := NewRatesController(mockService)

	t.Run("Success", func(t *testing.T) {
		expectedRate := &models.Rate{
			Ask:       1111.11,
			Bid:       22222.22,
			Timestamp: 3333333333,
		}

		mockService.EXPECT().
			GetCurrentRate(gomock.Any()).
			Return(expectedRate, nil).
			Times(1)

		resp, err := controller.GetRates(context.Background(), &pb.GetRatesRequest{})

		assert.NoError(t, err)
		assert.Equal(t, expectedRate.Ask, resp.Ask)
		assert.Equal(t, expectedRate.Bid, resp.Bid)
		assert.Equal(t, expectedRate.Timestamp, resp.Timestamp)
	})

	t.Run("Error", func(t *testing.T) {
		mockService.EXPECT().
			GetCurrentRate(gomock.Any()).
			Return(nil, assert.AnError).
			Times(1)

		_, err := controller.GetRates(context.Background(), &pb.GetRatesRequest{})
		assert.Error(t, err)
	})
}

func TestRatesController_HealthCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	controller := NewRatesController(nil)
	resp, err := controller.HealthCheck(context.Background(), &pb.HealthCheckRequest{})

	assert.NoError(t, err)
	assert.Equal(t, pb.HealthCheckResponse_SERVING, resp.Status)
}
