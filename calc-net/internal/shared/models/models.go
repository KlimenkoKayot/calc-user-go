package models

import (
	"time"

	"github.com/klimenkokayot/calc-net-go/internal/shared/customList"
)

const (
	StateDone       string = "Выполнено."
	StateInProgress string = "В обработке."
	StateError      string = "Ошибка."
)

// Пара токенов для валидации
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Структура арифметического выражения
type Expression struct {
	Id     string  `json:"id,omitempty"`
	Status string  `json:"status,omitempty"`
	Result float64 `json:"result,omitempty"`
	Value  string  `json:"expression,omitempty"`

	Hash [64]byte               `json:"-"`
	List *customList.LinkedList `json:"-"`
}

// Структура подзадачи выражения
type Task struct {
	Id             uint          `json:"id"`
	FirstArgument  float64       `json:"arg1"`
	SecondArgument float64       `json:"arg2"`
	Operation      rune          `json:"operation"`
	OperationTime  time.Duration `json:"operation_time"`
}

// Структура для обработки ответов на подзадачу
type TaskResult struct {
	Id     uint    `json:"id"`
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

// Структура для обработки запросов на статус запроса
type Result struct {
	Id     uint    `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

// Хранит хэш и флаг, если какой-то агент забрал доступ к запросу
type RequestExpression struct {
	Hash     [64]byte
	InAction bool
}
