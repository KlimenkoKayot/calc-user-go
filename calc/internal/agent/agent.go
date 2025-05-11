package agent

import (
	"log"

	config "github.com/klimenkokayot/calc-net-go/internal/agent/config"
	"github.com/klimenkokayot/calc-net-go/internal/agent/service"
)

// Структура агента, нужен порт и экземпляр сервиса
type Agent struct {
	Service          *service.AgentService
	OrchestratorPort int
}

// Создание нового агента
func NewAgent() (*Agent, error) {
	// Получение конфига
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	// Получение экземпляра агента
	service := service.NewAgentService(*config)
	return &Agent{
		service,
		config.OrchestratorPort,
	}, nil
}

// Запуск агента
func (a *Agent) Run() error {
	log.Println("start new agent")
	return a.Service.Run()
}
