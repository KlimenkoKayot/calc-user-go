package agent

import "fmt"

var (
	ErrNotFound       = fmt.Errorf("Endpoint не найден или нет задач")
	ErrInternalServer = fmt.Errorf("Ошибка на стороне оркестратора")
	ErrNotAuth        = fmt.Errorf("Доступ запрещен")
)
