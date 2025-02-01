package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppConf struct {
	AppName     string
	Port        string
	GarantexURL string
	ClickHouse  ClickHouseConfig
}

type ClickHouseConfig struct {
	Host     string
	Port     int
	DBName   string
	User     string
	Password string
}

func (a *AppConf) Init(logger *zap.SugaredLogger) {
	logger.Info("Initializing configuration...")

	// Автоматически загружает переменные окружения
	viper.AutomaticEnv()

	// Устанавливаем значения по умолчанию
	viper.SetDefault("APP_NAME", "getRates")

	// Список обязательных переменных окружения
	requiredVars := []string{
		"GARANTEX_API_URL",
		"CLICKHOUSE_HOST",
		"CLICKHOUSE_PORT",
		"CLICKHOUSE_DB",
		"CLICKHOUSE_USER",
		"CLICKHOUSE_PASSWORD",
	}

	// Проверяем все обязательные переменные
	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			logger.Fatalf("Required environment variable is missing: %s", key)
		}
	}

	a.AppName = viper.GetString("APP_NAME")
	a.Port = viper.GetString("PORT")
	a.GarantexURL = viper.GetString("GARANTEX_API_URL")
	a.ClickHouse.Host = viper.GetString("CLICKHOUSE_HOST")
	a.ClickHouse.Port = viper.GetInt("CLICKHOUSE_PORT")
	a.ClickHouse.DBName = viper.GetString("CLICKHOUSE_DB")
	a.ClickHouse.User = viper.GetString("CLICKHOUSE_USER")
	a.ClickHouse.Password = viper.GetString("CLICKHOUSE_PASSWORD")

	logger.Info("Configuration loaded successfully")
}
