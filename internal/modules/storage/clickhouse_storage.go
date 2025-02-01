package storage

import (
	"context"
	"fmt"
	"get-rates-usdt-grpc-service/config"
	"get-rates-usdt-grpc-service/internal/infrastracture/db"
	"get-rates-usdt-grpc-service/internal/models"
	"github.com/ClickHouse/clickhouse-go/v2"
	"time"
)

type RatesStorage interface {
	SaveRate(ctx context.Context, rate *models.Rate) error
}

type ClickHouseStorage struct {
	conn clickhouse.Conn
}

func NewClickHouseStorage(cfg config.ClickHouseConfig) (*ClickHouseStorage, error) {
	dns := fmt.Sprintf(
		"clickhouse://%s:%s@%s:%d/%s?x-multi-statement=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	if err := db.ApplyMigrations(dns); err != nil {
		return nil, fmt.Errorf("migrations failed: %w", err)
	}
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.DBName,
			Username: cfg.User,
			Password: cfg.Password,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	return &ClickHouseStorage{conn: conn}, nil
}

func (s *ClickHouseStorage) SaveRate(ctx context.Context, rate *models.Rate) error {
	return s.conn.Exec(ctx, `
		INSERT INTO rates (timestamp, ask, bid, created_at)
		VALUES (?, ?, ?, ?)`,
		rate.Timestamp,
		rate.Ask,
		rate.Bid,
		time.Now(),
	)
}
