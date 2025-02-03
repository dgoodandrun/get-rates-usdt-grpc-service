package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppConf struct {
	AppName     string
	Port        string
	MetricsPort string
	GarantexURL string
	Market      string
	Postgres    PostgresConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	DBName   string
	User     string
	Password string
}

//type ClickHouseConfig struct {
//	Host     string
//	Port     int
//	DBName   string
//	User     string
//	Password string
//}

func (a *AppConf) Init(logger *zap.SugaredLogger) {
	logger.Info("Initializing configuration...")

	// Автоматически загружает переменные окружения
	viper.AutomaticEnv()

	// Устанавливаем значения по умолчанию
	viper.SetDefault("APP_NAME", "getRates")

	// Список обязательных переменных окружения
	requiredVars := []string{
		"PORT",
		"METRICS_PORT",
		"GARANTEX_API_URL",
		"GARANTEX_API_URL_MARKET",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
	}

	// Проверяем все обязательные переменные
	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			logger.Fatalf("Required environment variable is missing: %s", key)
		}
	}

	a.AppName = viper.GetString("APP_NAME")
	a.Port = viper.GetString("PORT")
	a.MetricsPort = viper.GetString("METRICS_PORT")
	a.GarantexURL = viper.GetString("GARANTEX_API_URL")
	a.Market = viper.GetString("GARANTEX_API_URL_MARKET")
	a.Postgres.Host = viper.GetString("POSTGRES_HOST")
	a.Postgres.Port = viper.GetInt("POSTGRES_PORT")
	a.Postgres.DBName = viper.GetString("POSTGRES_DB")
	a.Postgres.User = viper.GetString("POSTGRES_USER")
	a.Postgres.Password = viper.GetString("POSTGRES_PASSWORD")
	//a.ClickHouse.Host = viper.GetString("CLICKHOUSE_HOST")
	//a.ClickHouse.Port = viper.GetInt("CLICKHOUSE_PORT")
	//a.ClickHouse.DBName = viper.GetString("CLICKHOUSE_DB")
	//a.ClickHouse.User = viper.GetString("CLICKHOUSE_USER")
	//a.ClickHouse.Password = viper.GetString("CLICKHOUSE_PASSWORD")

	logger.Info("Configuration loaded successfully")
}
