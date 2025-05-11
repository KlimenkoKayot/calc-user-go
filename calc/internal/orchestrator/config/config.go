package orchestrator

import (
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
)

// Конфиг для работы оркестратора
type Config struct {
	Port                  int    // порт для запуска сервера
	TimeAdditionMs        uint64 // время выполнения операции сложения в миллисекундах
	TimeSubtractionMs     uint64 // время выполнения операции вычитания в миллисекундах
	TimeMultiplicationsMs uint64 // время выполнения операции умножения в миллисекундах
	TimeDivisionsMs       uint64 // время выполнения операции деления в миллисекундах
	AuthBaseURL           string
}

// Загрузка .env из любой точки входа в проекте
// (нужно для запуска тестов)
func init() {
	// Загрузка переменных окружения из .env файла
	projectDirName := `calc-net-go`
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	if err := godotenv.Load(string(rootPath) + `/.env`); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// Создает новый конфиг, получает значения переменных из .env
func NewConfig() (*Config, error) {
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

	TimeAdditionString := os.Getenv("TIME_ADDITION_MS")
	if TimeAdditionString == "" {
		TimeAdditionString = "0"
	}
	TimeAdditionMs, err := strconv.Atoi(TimeAdditionString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if TimeAdditionMs < 0 {
		return nil, ErrInvalidTime
	}

	TimeSubtractionString := os.Getenv("TIME_SUBTRACTION_MS")
	if TimeSubtractionString == "" {
		TimeSubtractionString = "0"
	}
	TimeSubtractionMs, err := strconv.Atoi(TimeSubtractionString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if TimeSubtractionMs < 0 {
		return nil, ErrInvalidTime
	}

	TimeMultiplicationsString := os.Getenv("TIME_MULTIPLICATIONS_MS")
	if TimeMultiplicationsString == "" {
		TimeMultiplicationsString = "0"
	}
	TimeMultiplicationsMs, err := strconv.Atoi(TimeMultiplicationsString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if TimeMultiplicationsMs < 0 {
		return nil, ErrInvalidTime
	}

	TimeDivisionsString := os.Getenv("TIME_DIVISIONS_MS")
	if TimeDivisionsString == "" {
		TimeDivisionsString = "0"
	}
	TimeDivisionsMs, err := strconv.Atoi(TimeDivisionsString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if TimeDivisionsMs < 0 {
		return nil, ErrInvalidTime
	}

	authBaseURL := os.Getenv("AUTH_BASE_URL")

	return &Config{
		Port,
		uint64(TimeAdditionMs),
		uint64(TimeSubtractionMs),
		uint64(TimeMultiplicationsMs),
		uint64(TimeDivisionsMs),
		authBaseURL,
	}, nil
}
