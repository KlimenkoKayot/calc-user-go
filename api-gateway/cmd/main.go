package main

import (
	"fmt"
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

	fmt.Println(cfg.ApiGateway.Logger)

	logger, err := logger.NewAdapter(&logger.Config{
		Adapter: cfg.ApiGateway.Logger,
		Level:   logger.LevelDebug,
	})
	if err != nil {
		log.Fatal(err)
	}

	loggerService := logger.WithLayer("SERVICE")
	service, err := service.NewProxyService(loggerService, cfg)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	loggerHandler := logger.WithLayer("HANDLER")
	handler, err := handler.NewProxyHandler(service, loggerHandler, cfg)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	loggerServer := logger.WithLayer("SERVER")
	server, err := server.NewProxyServer(handler, loggerServer, cfg)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	loggerApp := logger.WithLayer("APP")
	app, err := app.NewProxyApplication(server, loggerApp, cfg)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	if err := app.Run(); err != nil {
		logger.Fatal(err.Error())
	}
}
