package agent

import (
	"log"
	"net/http"
	"time"

	transport "github.com/klimenkokayot/calc-net-go/internal/agent/worker/transport"
	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
)

// Воркер решает задачи, которые получает от оркестратора
type Worker struct {
	Client    *http.Client  // Клиент для выполнения запросов
	URL       string        // Адрес оркестратора для запросов
	SleepTime time.Duration // Задержка между запросами агента
}

// Создаем нового воркера, нужен адрес оркестратора и задержка из .env
func NewWorker(url string, sleepTime time.Duration) *Worker {
	return &Worker{
		&http.Client{},
		url,
		sleepTime,
	}
}

// Решение полученной задачи
func (w *Worker) Solve(task *models.Task) *models.TaskResult {
	log.Printf("Worker получил новую подзадачу c id: %d\n", task.Id)
	time.Sleep(time.Millisecond * task.OperationTime)
	var answer float64
	var err error
	switch task.Operation {
	case '+':
		answer = task.FirstArgument + task.SecondArgument
	case '-':
		answer = task.FirstArgument - task.SecondArgument
	case '*':
		answer = task.FirstArgument * task.SecondArgument
	case '/':
		if task.SecondArgument == 0 {
			err = ErrDivisionByZero
			answer = 0
			break
		}
		answer = task.FirstArgument / task.SecondArgument
	}
	log.Printf("Получен ответ на подзадачу с id: %d, ответ: %f\n", task.Id, answer)
	errString := ""
	if err != nil {
		errString = err.Error()
	}
	return &models.TaskResult{
		Id:     task.Id,
		Result: answer,
		Error:  errString,
	}
}

// Попытка получения подзадачи и ее решения
func (w *Worker) Process() error {
	task, err := transport.GetTask(w.Client, w.URL)
	if err != nil {
		return err
	}
	result := w.Solve(task)
	err = transport.PostTask(w.Client, w.URL, result)
	return err
}

// Воркер начинает работу, GET -> SOLVE -> POST -> SLEEP
func (w *Worker) Run() error {
	for {
		err := w.Process()
		if err != nil {
			log.Println(err.Error())
		}
		time.Sleep(w.SleepTime)
	}
}
