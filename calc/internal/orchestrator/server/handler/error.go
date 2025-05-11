package orchestrator

import "fmt"

var (
	ErrInternalServer    = fmt.Errorf("у нас что-то пошло не так")
	ErrInvalidBodyDecode = fmt.Errorf("не удалось прочитать тело запроса")
)
