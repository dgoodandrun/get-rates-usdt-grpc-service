package storage

import (
	"context"
	"database/sql"
	"fmt"
	"get-rates-usdt-grpc-service/config"
	"get-rates-usdt-grpc-service/internal/infrastracture/db"
	"get-rates-usdt-grpc-service/internal/infrastracture/metrics"
	"get-rates-usdt-grpc-service/internal/models"
	"time"
)

//go:generate mockgen -source=postgres_storage.go -destination=mocks/postgres_storage_mock.go -package=mocks

type RatesStorage interface {
	SaveRate(ctx context.Context, rate *models.Rate) error
	HealthCheck(ctx context.Context) error
}

type postgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(cfg config.PostgresConfig) (RatesStorage, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	// Применяем миграции
	if err := db.ApplyMigrations(dsn); err != nil {
		return nil, fmt.Errorf("migrations failed: %w", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &postgresStorage{db: db}, nil
}

func (s *postgresStorage) SaveRate(ctx context.Context, rate *models.Rate) error {
	start := time.Now()

	_, err := s.db.ExecContext(ctx, `
        INSERT INTO rates (timestamp, ask, bid, created_at)
        VALUES ($1, $2, $3, $4)`,
		rate.Timestamp,
		rate.Ask,
		rate.Bid,
		time.Now(),
	)

	duration := time.Since(start).Seconds()
	operation := "insert"
	if err != nil {
		operation = "insert_error"
	}

	metrics.DBRequests.WithLabelValues(operation).Inc()
	metrics.DBDuration.WithLabelValues(operation).Observe(duration)

	return err
}

func (s *postgresStorage) HealthCheck(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
