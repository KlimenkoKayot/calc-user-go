package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	worker "github.com/klimenkokayot/calc-net-go/internal/agent/worker"
	"github.com/klimenkokayot/calc-user-go/config"
)

// Методы бизнес-логики агента
type AgentService struct {
	OrchestratorUrl string        // Адрес оркестратора
	AgentSleepTime  time.Duration // Задержка между запросами в оркестратор
	ComputingPower  uint64        // Кол-во горутин, которые запускаются внутри агента
	wg              *sync.WaitGroup
}

// Новый агент, конфиг обязателен
func NewAgentService(config config.Config) *AgentService {
	return &AgentService{
		fmt.Sprintf("http://127.0.0.1:%d/internal/task", config.Calc.Orchestrator.Port),
		config.Calc.Agent.Timeout,
		uint64(config.Calc.Agent.Workers),
		&sync.WaitGroup{},
	}
}

// Запуск нового воркера
func (s *AgentService) StartNewWorker() error {
	log.Println("Запущен новый worker...")
	s.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		worker := worker.NewWorker(s.OrchestratorUrl, s.AgentSleepTime)
		worker.Run()
		wg.Done()
	}(s.wg)
	return nil
}

// Старт сервиса, запускает воркеров в отдельных горутинах
func (s *AgentService) Run() error {
	for i := 0; i < int(s.ComputingPower); i++ {
		s.StartNewWorker()
	}
	s.wg.Wait()
	return nil
}
