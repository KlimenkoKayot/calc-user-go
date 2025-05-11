package orchestrator_test

import (
	"log"
	"testing"
	"time"

	service "github.com/klimenkokayot/calc-net-go/internal/orchestrator/service"
	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
	"github.com/klimenkokayot/calc-user-go/config"
)

type TestCase struct {
	Value    interface{}
	Expected interface{}
	IsError  bool
}

func TestService(t *testing.T) {
	// Запуск тестов для каждого метода
	config, err := config.Load(config.GetConfigPath())
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	t.Run("NewExpression", func(t *testing.T) {
		testServiceNewExpression(t, config)
	})
	t.Run("AddExpression", func(t *testing.T) {
		testServiceAddExpression(t, config)
	})
	t.Run("GetAllExpressions", func(t *testing.T) {
		testServiceGetAllExpressions(t, config)
	})
	t.Run("FindNewTasks", func(t *testing.T) {
		testServiceFindNewTasks(t, config)
	})
	t.Run("GetTask", func(t *testing.T) {
		testServiceGetTask(t, config)
	})
	t.Run("ProcessErrorAnswer", func(t *testing.T) {
		testServiceProcessErrorAnswer(t, config)
	})
	t.Run("ProcessAnswer", func(t *testing.T) {
		testServiceProcessAnswer(t, config)
	})
	t.Run("OperationTime", func(t *testing.T) {
		testServiceOperationTime(t, config)
	})
}

func testServiceNewExpression(t *testing.T, config *config.Config) {
	service := service.NewOrchestratorService(config)
	testCases := []TestCase{
		{
			Value:    "2+2*2",
			Expected: nil,
			IsError:  false,
		},
		{
			Value:    "2",
			Expected: nil,
			IsError:  false,
		},
		{
			Value:    "",
			Expected: nil,
			IsError:  true,
		},
		{
			Value:    "2+*2",
			Expected: nil,
			IsError:  true,
		},
		{
			Value:    "2+(3+5)",
			Expected: nil,
			IsError:  false,
		},
		// Выражение задано корректно, но упадет на этапе подсчета
		{
			Value:    "1/0",
			Expected: nil,
			IsError:  false,
		},
		{
			Value:    "1 o 0",
			Expected: nil,
			IsError:  true,
		},
	}
	for _, test := range testCases {
		_, err := service.NewExpression(test.Value.(string))
		if err != nil && !test.IsError {
			t.Errorf("want not err, got error: %s\nexpression: %s\n", err.Error(), test.Value.(string))
		}
		if err == nil && test.IsError {
			t.Errorf("want err, got nil:\nexpression: %s\n", test.Value.(string))
		}
	}
}

func testServiceAddExpression(t *testing.T, config *config.Config) {
	service := service.NewOrchestratorService(config)
	testCases := []TestCase{
		{
			Value:    "2+2",
			Expected: nil,
			IsError:  false,
		},
		{
			Value:    "",
			Expected: nil,
			IsError:  true,
		},
		{
			Value:    "2 o 9",
			Expected: nil,
			IsError:  true,
		},
	}

	for _, test := range testCases {
		_, err := service.AddExpression(test.Value.(string))
		if err != nil && !test.IsError {
			t.Errorf("want no error, got error: %s\nexpression: %s\n", err.Error(), test.Value.(string))
		}
		if err == nil && test.IsError {
			t.Errorf("want error, got nil:\nexpression: %s\n", test.Value.(string))
		}
	}
}

func testServiceGetAllExpressions(t *testing.T, config *config.Config) {
	service := service.NewOrchestratorService(config)

	// Добавляем выражения
	_, _ = service.AddExpression("2+2")
	_, _ = service.AddExpression("3*3")

	// Добавляем ответы
	service.Answers[[64]byte{}] = 3

	// Получаем все выражения
	expressions := service.GetAllExpressions()

	// Проверяем количество выражений
	if len(expressions) != 3 {
		t.Errorf("want 3 expressions, got %d", len(expressions))
	}
}

func testServiceFindNewTasks(t *testing.T, config *config.Config) {
	service := service.NewOrchestratorService(config)

	// Добавляем выражение
	expression, _ := service.NewExpression("2+2*2")

	// Ищем подзадачи
	tasks, err := service.FindNewTasks(expression)
	if err != nil {
		t.Errorf("want no error, got error: %s", err.Error())
	}

	// Проверяем количество подзадач
	if len(tasks) != 1 {
		t.Errorf("want 1 task, got %d", len(tasks))
	}
}

func testServiceGetTask(t *testing.T, config *config.Config) {
	service := service.NewOrchestratorService(config)

	// Добавляем выражение
	_, _ = service.AddExpression("2+2*2")

	// Получаем подзадачу
	task, err := service.GetTask()
	if err != nil {
		t.Errorf("want no error, got error: %s", err.Error())
	}

	// Проверяем, что подзадача не пустая
	if task == nil {
		t.Error("want task, got nil")
	}
}

func testServiceProcessErrorAnswer(t *testing.T, config *config.Config) {
	service := service.NewOrchestratorService(config)

	// Добавляем выражение
	_, _ = service.AddExpression("2+2/0")

	// Получаем подзадачу
	task, _ := service.GetTask()

	// Обрабатываем ошибку
	service.ProcessErrorAnswer(&models.TaskResult{
		Id:    task.Id,
		Error: "division by zero",
	})

	// Проверяем, что выражение помечено как ошибочное
	expressions := service.GetAllExpressions()
	if expressions[0].Status != models.StateError {
		t.Error("want expression status to be error, got", expressions[0].Status)
	}
}

func testServiceProcessAnswer(t *testing.T, config *config.Config) {
	service := service.NewOrchestratorService(config)

	// Добавляем выражение
	_, _ = service.AddExpression("2+2*2")

	// Получаем подзадачу
	task, _ := service.GetTask()

	// Обрабатываем ответ
	service.ProcessAnswer(&models.TaskResult{
		Id:     task.Id,
		Result: 4.0,
	})

	// Проверяем, что выражение обновлено
	expressions := service.GetAllExpressions()
	if len(expressions) != 1 {
		t.Errorf("want 1 expression, got %d", len(expressions))
	}
}

func testServiceOperationTime(t *testing.T, config *config.Config) {
	service := service.NewOrchestratorService(config)

	testCases := []TestCase{
		{
			Value:    '+',
			Expected: service.TimeAdditionMs,
			IsError:  false,
		},
		{
			Value:    '-',
			Expected: service.TimeSubtractionMs,
			IsError:  false,
		},
		{
			Value:    '*',
			Expected: service.TimeMultiplicationsMs,
			IsError:  false,
		},
		{
			Value:    '/',
			Expected: service.TimeDivisionsMs,
			IsError:  false,
		},
		{
			Value: '^',
			// Типа нулевой time.Duration
			Expected: 0 * time.Second,
			IsError:  true,
		},
	}

	for _, test := range testCases {
		duration, err := service.OperationTime(test.Value.(rune))
		if err != nil && !test.IsError {
			t.Errorf("want no error, got error: %s\noperation: %c\n", err.Error(), test.Value.(rune))
		} else if err == nil && test.IsError {
			t.Errorf("want error, got nil:\noperation: %c\n", test.Value.(rune))
		} else if duration != test.Expected.(time.Duration) {
			t.Errorf("want duration %v, got %v\noperation: %c\n", test.Expected, duration, test.Value.(rune))
		}
	}
}
