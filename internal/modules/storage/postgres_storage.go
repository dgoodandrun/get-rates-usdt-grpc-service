package storage

import (
	"context"
	"fmt"
	"get-rates-usdt-grpc-service/config"
	"get-rates-usdt-grpc-service/internal/infrastracture/db"
	"get-rates-usdt-grpc-service/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

//go:generate mockgen -source=postgres_storage.go -destination=mocks/postgres_storage_mock.go -package=mocks

type RatesStorage interface {
	SaveRate(ctx context.Context, rate *models.Rate) error
}

type PGXPool interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Close()
}

type PostgresStorage struct {
	pool PGXPool
}

func newPostgresStorageWithPool(pool PGXPool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func NewPostgresStorage(cfg config.PostgresConfig) (*PostgresStorage, error) {
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

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return newPostgresStorageWithPool(pool), nil
}

func (s *PostgresStorage) SaveRate(ctx context.Context, rate *models.Rate) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO rates (timestamp, ask, bid, created_at)
		VALUES ($1, $2, $3, $4)`,
		rate.Timestamp,
		rate.Ask,
		rate.Bid,
		time.Now(),
	)

	return err
}
