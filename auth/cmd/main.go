package main

import (
	"log"

	"github.com/klimenkokayot/avito-go/libs/logger"
	"github.com/klimenkokayot/avito-go/services/auth/config"
	"github.com/klimenkokayot/avito-go/services/auth/internal/app"
)

func main() {
	config, err := config.Load("")
	if err != nil {
		log.Fatalf("Ошибка при инициализации config`a: %s.", err.Error())
	}

	logger, err := logger.NewAdapter(&logger.Config{
		Adapter: config.Logger,
		Level:   logger.LevelDebug - 1,
	})
	if err != nil {
		log.Fatalf("Ошибка при инициализации config`a: %s.", err.Error())
	}

	appLogger := logger.WithLayer("APP")
	app, err := app.NewApplication(config, appLogger)
	if err != nil {
		appLogger.Fatal("Ошибка при инициализации application: " + err.Error() + ".")
	}
	if err := app.Run(); err != nil {
		appLogger.Fatal(err.Error())
	}
}
