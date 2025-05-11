package utils

import "fmt"

var (
	ErrInvalidBase64Decode = fmt.Errorf("неверный формат для хэш-функции (SHA512 должен иметь 64 байта)")
)
