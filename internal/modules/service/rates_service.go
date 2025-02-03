package service

import (
	"context"
	"encoding/json"
	"fmt"
	"get-rates-usdt-grpc-service/internal/models"
	"get-rates-usdt-grpc-service/internal/modules/storage"
	"io"
	"net/http"
)

type RatesService interface {
	GetCurrentRate(ctx context.Context) (*models.Rate, error)
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
	resp, err := http.Get(fmt.Sprintf(s.apiURL, s.market))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
		fmt.Sscanf(data.Asks[0].Price, "%f", &ask)
	}
	if len(data.Bids) > 0 {
		fmt.Sscanf(data.Bids[0].Price, "%f", &bid)
	}

	rate := &models.Rate{
		Ask:       ask,
		Bid:       bid,
		Timestamp: data.Timestamp,
	}

	if err := s.storage.SaveRate(ctx, rate); err != nil {
		return rate, fmt.Errorf("failed to save rate: %w", err)
	}

	return rate, nil
}
