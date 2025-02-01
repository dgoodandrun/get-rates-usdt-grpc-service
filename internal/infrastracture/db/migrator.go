package db

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigrations() error {
	m, err := migrate.New(
		"file://migrations",
		"clickhouse://user:pass@host:port/dbname?x-multi-statement=true",
	)
	if err != nil {
		return err
	}
	return m.Up()
}
