package storage

import (
	"context"

	"testing"

	"get-rates-usdt-grpc-service/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPostgresStorage_SaveRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Error creating mock database:", err)
	}
	defer db.Close()

	storage := &postgresStorage{db: db}
	testRate := &models.Rate{
		Ask:       1111.11,
		Bid:       2222.22,
		Timestamp: 33333333,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO rates").
			WithArgs(
				testRate.Timestamp,
				testRate.Ask,
				testRate.Bid,
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := storage.SaveRate(context.Background(), testRate)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO rates").
			WithArgs(
				testRate.Timestamp,
				testRate.Ask,
				testRate.Bid,
				sqlmock.AnyArg(),
			).
			WillReturnError(assert.AnError)

		err := storage.SaveRate(context.Background(), testRate)
		assert.Error(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
