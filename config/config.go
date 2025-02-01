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
	viper.SetDefault("APP_NAME", "geo")

	// Список обязательных переменных окружения
	requiredVars := []string{
		"API_KEY",
		"SECRET_KEY",
		"GEO_SERVICE_PORT",
		"REDIS_ADDR",
		"RATE_LIMIT_PER_MINUTE",
		"QUEUE",
		"TOPIC",
		"RABBITMQ_ADDR",
		"KAFKA_ADDR",
	}

	// Проверяем все обязательные переменные
	for _, key := range requiredVars {
		if !viper.IsSet(key) {
			logger.Fatalf("Required environment variable is missing: %s", key)
		}
	}

	a.AppName = viper.GetString("APP_NAME")
	a.APIKey = viper.GetString("API_KEY")
	a.SecretKey = viper.GetString("SECRET_KEY")
	a.Port = viper.GetString("GEO_SERVICE_PORT")
	a.RedisAddr = viper.GetString("REDIS_ADDR")
	a.Rate = viper.GetInt("RATE_LIMIT_PER_MINUTE")
	a.Queue = viper.GetString("QUEUE")
	a.Topic = viper.GetString("TOPIC")
	a.RabbitMQAddr = viper.GetString("RABBITMQ_ADDR")
	a.KafkaAddr = viper.GetString("KAFKA_ADDR")

	logger.Info("Configuration loaded successfully")
}
