package storage

import (
	"context"
	"errors"
	"get-rates-usdt-grpc-service/internal/models"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgresStorage_SaveRate(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal("Error creating mock pool:", err)
	}
	defer mockPool.Close()

	storage := &PostgresStorage{pool: mockPool}
	testRate := &models.Rate{
		Ask:       1111.11,
		Bid:       2222.22,
		Timestamp: 33333333,
	}

	t.Run("Success", func(t *testing.T) {
		mockPool.ExpectExec("INSERT INTO rates").
			WithArgs(
				testRate.Timestamp,
				testRate.Ask,
				testRate.Bid,
				pgxmock.AnyArg(),
			).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := storage.SaveRate(context.Background(), testRate)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockPool.ExpectExec("INSERT INTO rates").
			WithArgs(
				testRate.Timestamp,
				testRate.Ask,
				testRate.Bid,
				pgxmock.AnyArg(),
			).
			WillReturnError(errors.New("database error"))

		err := storage.SaveRate(context.Background(), testRate)
		assert.ErrorContains(t, err, "database error")
	})

	assert.NoError(t, mockPool.ExpectationsWereMet())
}
