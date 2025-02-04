package service

import (
	"context"
	"encoding/json"
	"fmt"
	"get-rates-usdt-grpc-service/internal/infrastructure/errors"
	"get-rates-usdt-grpc-service/internal/infrastructure/metrics"
	"get-rates-usdt-grpc-service/internal/models"
	"get-rates-usdt-grpc-service/internal/modules/storage"
	"io"
	"net/http"
	"strconv"
	"time"
)

//go:generate mockgen -source=rates_service.go -destination=mocks/rates_service_mock.go -package=mocks

type RatesService interface {
	GetCurrentRate(ctx context.Context) (*models.Rate, error)
	HealthCheck(ctx context.Context) error
}

type ratesService struct {
	storage storage.RatesStorage
	apiURL  string
	market  string
}

func NewRatesService(storage storage.RatesStorage, apiURL string, market string) RatesService {
	return &ratesService{
		storage: storage,
		apiURL:  apiURL,
		market:  market,
	}
}

func (s *ratesService) GetCurrentRate(ctx context.Context) (*models.Rate, error) {
	start := time.Now()
	resp, err := http.Get(fmt.Sprintf(s.apiURL, s.market))
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil || resp.StatusCode != http.StatusOK {
		status = "error"
		metrics.ExternalAPIRequests.WithLabelValues(status).Observe(duration)
		metrics.ExternalAPIDuration.Observe(duration)
		return nil, err
	}
	defer resp.Body.Close()

	metrics.ExternalAPIRequests.WithLabelValues(status).Observe(duration)
	metrics.ExternalAPIDuration.Observe(duration)

	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Asks []struct {
			Price string `json:"price"`
		} `json:"asks"`
		Bids []struct {
			Price string `json:"price"`
		} `json:"bids"`
		Timestamp int64 `json:"timestamp"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	ask, bid := 0.0, 0.0
	if len(data.Asks) > 0 {
		ask, _ = strconv.ParseFloat(data.Asks[0].Price, 64)
	}
	if len(data.Bids) > 0 {
		bid, _ = strconv.ParseFloat(data.Bids[0].Price, 64)
	}

	rate := &models.Rate{
		Ask:       ask,
		Bid:       bid,
		Timestamp: data.Timestamp,
	}

	if err := s.storage.SaveRate(ctx, rate); err != nil {
		return rate, errors.SaveRateError
	}

	return rate, nil
}

func (s *ratesService) HealthCheck(ctx context.Context) error {
	resp, err := http.Get(fmt.Sprintf(s.apiURL, s.market))
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.ApiError
	}
	defer resp.Body.Close()

	if err := s.storage.HealthCheck(ctx); err != nil {
		return errors.StorageError
	}

	return nil
}
