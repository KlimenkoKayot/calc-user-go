package rpn

import (
	"errors"
)

var (
	ErrInvalidExpression       = errors.New("invalid expression")
	ErrInvalidSymbolExpression = errors.New("invalid symbol in expression")
	ErrDevisionByZero          = errors.New("division by zero")
	ErrEmptyExpression         = errors.New("empty expression")
)
