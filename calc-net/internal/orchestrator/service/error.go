package orchestrator

import "fmt"

var (
	ErrInvalidSymbolRPN  = fmt.Errorf("неизвестный символ в библиотеке rpn")
	ErrInvalidOperation  = fmt.Errorf("неизвестная операция при попытке поиска времени выполнения")
	ErrHaveNoTask        = fmt.Errorf("нет задач")
	ErrAnswerExpression  = fmt.Errorf("попытка поиска задач в выражении, состоящем из ответа")
	ErrInvalidExpression = fmt.Errorf("выражение задано некорректно")
	ErrZeroExpression    = fmt.Errorf("пустое выражение")
)
