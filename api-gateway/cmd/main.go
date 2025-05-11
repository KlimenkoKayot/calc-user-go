package main

import (
	"log"

	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/app"
	server "github.com/klimenkokayot/calc-user-go/api-gateway/internal/infrastructure/http"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/infrastructure/service"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/interfaces/handler"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/logger"
	"github.com/klimenkokayot/calc-user-go/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.NewAdapter(&logger.Config{
		Adapter: cfg.ApiGateway.Logger,
		Level:   logger.LevelDebug,
	})
	if err != nil {
		log.Fatal(err)
	}

	service, err := service.NewProxyService(logger, cfg)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	handler, err := handler.NewProxyService(service, logger, cfg)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	server, err := server.NewProxyServer(handler, logger, cfg)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	app, err := app.NewProxyApplication(server, logger, cfg)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	if err := app.Run(); err != nil {
		logger.Fatal(err.Error())
	}
}
