package main

import (
	"log"

	client "github.com/klimenkokayot/calc-net-go/internal/orchestrator/client"
	orchestrator "github.com/klimenkokayot/calc-net-go/internal/orchestrator/server"
	handler "github.com/klimenkokayot/calc-net-go/internal/orchestrator/server/handler"
	service "github.com/klimenkokayot/calc-net-go/internal/orchestrator/service"
	"github.com/klimenkokayot/calc-net-go/pkg/logger"
	"github.com/klimenkokayot/calc-user-go/config"
)

func main() {
	// Создаем конфиг
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	// Создаем новый сервер aka оркестратора
	orchestratorService := service.NewOrchestratorService(config)
	orchestratorHandler := handler.NewOrchestratorHandler(config, orchestratorService)
	logger, err := logger.NewAdapter(&logger.Config{
		Adapter: logger.AdapterZap,
		Level:   logger.LevelDebug,
	})
	if err != nil {
		log.Fatal(err)
	}
	orchestratorClient := client.NewAuthClient(config.ApiGateway.Services.Auth.URL, logger)
	server, err := orchestrator.NewServer(config, orchestratorHandler, orchestratorClient)
	if err != nil {
		log.Fatal(err)
	}
	// Запуск созданного сервера
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
