package service

import (
	"context"
	"get-rates-usdt-grpc-service/internal/modules/storage/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRatesService_GetCurrentRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockRatesStorage(ctrl)
	mockStorage.EXPECT().SaveRate(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
            "asks": [{"price": "111.11"}],
            "bids": [{"price": "222.22"}],
            "timestamp": 333333333
        }`))
	}))
	defer testServer.Close()

	t.Run("Success", func(t *testing.T) {
		rateService := NewRatesService(mockStorage, testServer.URL+"/?market=%s", "usdt")
		rate, err := rateService.GetCurrentRate(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 111.11, rate.Ask)
		assert.Equal(t, 222.22, rate.Bid)
		assert.Equal(t, int64(333333333), rate.Timestamp)
	})

	t.Run("HTTP Error", func(t *testing.T) {
		brokenService := NewRatesService(mockStorage, "invalid-url", "usdt")
		_, err := brokenService.GetCurrentRate(context.Background())
		assert.Error(t, err)
	})
}
