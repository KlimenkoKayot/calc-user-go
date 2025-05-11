package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Login     string    `db:"login" json:"login"`
	Secret    string    `db:"secret" json:"pass"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

var (
	ErrBadPassword  = fmt.Errorf("неправильный пароль")
	ErrUserNotFound = fmt.Errorf("пользователь не найден")
	ErrUserExist    = fmt.Errorf("пользователь уже существует")
)

type UserRepository interface {
	Add(login string, secret string) error
	Check(login string, pass string) (bool, error)
	FindByLogin(login string) (*User, error)
	ExistByLogin(login string) (bool, error)
}
