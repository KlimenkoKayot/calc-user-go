package agent

import "fmt"

var (
	ErrLoadEnvironment       = fmt.Errorf("ошибка загрузки переменных среды")
	ErrInvalidVariableType   = fmt.Errorf("неверный тип переменной среды")
	ErrInvalidPort           = fmt.Errorf("некорректный порт")
	ErrInvalidComputingValue = fmt.Errorf("число агентов (горутин) должно быть больше 0")
	ErrInvalidTime           = fmt.Errorf("время не может быть отрицательным значением")
)
