package agent

import (
	"log"

	"github.com/klimenkokayot/calc-net-go/internal/agent/service"
	"github.com/klimenkokayot/calc-user-go/config"
)

// Структура агента, нужен порт и экземпляр сервиса
type Agent struct {
	Service *service.AgentService
}

// Создание нового агента
func NewAgent() (*Agent, error) {
	// Получение конфига
	config, err := config.Load()
	if err != nil {
		return nil, err
	}
	// Получение экземпляра агента
	service := service.NewAgentService(*config)
	return &Agent{
		service,
	}, nil
}

// Запуск агента
func (a *Agent) Run() error {
	log.Println("start new agent")
	return a.Service.Run()
}
