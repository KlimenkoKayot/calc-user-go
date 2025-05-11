package agent

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	OrchestratorPort int           // Порт, на котором запускается оркестратор
	ComputingPower   uint64        // Количество горутин внутри агента
	AgentSleepTime   time.Duration // Задержка между запросами агента в оркестратор
}

// Создание нового конфига
func NewConfig() (*Config, error) {
	// Загрузка переменных среды
	if err := godotenv.Load(); err != nil {
		return nil, ErrLoadEnvironment
	}

	PortString := os.Getenv("PORT")
	if PortString == "" {
		PortString = "8080"
	}
	Port, err := strconv.Atoi(PortString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if Port < 0 {
		return nil, ErrInvalidPort
	}

	ComputingPowerString := os.Getenv("COMPUTING_POWER")
	if ComputingPowerString == "" {
		ComputingPowerString = "4"
	}
	ComputingPower, err := strconv.Atoi(ComputingPowerString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if ComputingPower < 0 {
		return nil, ErrInvalidComputingValue
	}

	sleepTimeString := os.Getenv("AGENT_SLEEP_TIME")
	if ComputingPowerString == "" {
		ComputingPowerString = "100"
	}
	sleepTime, err := strconv.Atoi(sleepTimeString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if sleepTime < 0 {
		return nil, ErrInvalidTime
	}

	return &Config{
		Port,
		uint64(ComputingPower),
		time.Duration(time.Millisecond * time.Duration(sleepTime)),
	}, nil
}
