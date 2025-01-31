package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

func NewLogger() *zap.SugaredLogger {
	// Устанавливаем конфигурацию для кастомного логгера
	config := zap.NewProductionConfig()
	config.Encoding = "console" // Для текстового вывода вместо JSON
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	// Настраиваем кастомный форматтер с цветами и переносами строк
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Устанавливает цвет в зависимости от уровня
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	// Создаем логгер с данной конфигурацией
	logger, err := config.Build()
	if err != nil {
		log.Fatal("failed to initialize logger: ", err)
	}

	return logger.Sugar()
}
