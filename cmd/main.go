package main

import (
	"os"
)

func main() {
	logger := logs.NewLogger()
	defer logger.Sync()

	conf := config.AppConf{}
	conf.Init(logger)

	app := run.NewApp(conf, logger)

	exitCode := app.
		// Инициализируем приложение
		Bootstrap().
		// Запускаем приложение
		Run()

	os.Exit(exitCode)
}
