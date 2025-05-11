package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/klimenkokayot/avito-go/libs/logger"
	domain "github.com/klimenkokayot/avito-go/services/auth/internal/domain/repository"
	"github.com/klimenkokayot/calc-user-go/config"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db     *sqlx.DB
	logger logger.Logger
	cfg    *config.Config
}

func NewUserRepository(cfg *config.Config, repoLogger logger.Logger) (domain.UserRepository, error) {
	repoLogger.Info("Инициализация user-репозитория.")
	dsn := cfg.Auth.Database.DSN
	if dsn == "" {
		// Установим значение по умолчанию для SQLite
		dsn = "file:auth.db?cache=shared&mode=rwc"
	}

	repoLogger.Info("Подключение по DSN.")
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		repoLogger.Error("Ошибка при подключении к sqlx.", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		return nil, fmt.Errorf("ошибка при подключении к sqlx: %w", err)
	}
	repoLogger.OK("Подключение к базе данных выполнено.")

	// Включение foreign key constraints для SQLite
	_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
	if err != nil {
		repoLogger.Error("Ошибка при включении foreign keys.", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		return nil, fmt.Errorf("ошибка при включении foreign keys: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			login TEXT UNIQUE NOT NULL,
			secret TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		repoLogger.Error("Ошибка при создании таблицы.", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		return nil, fmt.Errorf("ошибка при создании таблицы: %w", err)
	}

	repoLogger.OK("Успешно.")
	return &UserRepository{
		db,
		repoLogger,
		cfg,
	}, nil
}

func (ur *UserRepository) FindByLogin(login string) (*domain.User, error) {
	user := &domain.User{}
	err := ur.db.Get(
		user,
		"SELECT id, login, secret, created_at FROM users WHERE login = ?",
		login,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске по логину: %w", err)
	}
	return user, nil
}

func (ur *UserRepository) ExistByLogin(login string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE login = ?)"
	err := ur.db.QueryRow(query, login).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при exist по логину: %w", err)
	}
	return exists, nil
}

func (ur *UserRepository) Add(login string, secret string) error {
	found, err := ur.ExistByLogin(login)
	if err != nil {
		return fmt.Errorf("ошибка при проверки на существование: %w", err)
	} else if found {
		return domain.ErrUserExist
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("ошибка при генерации UUID: %w", err)
	}

	id := uuid.String()
	_, err = ur.db.Exec(
		"INSERT INTO users (id, login, secret) VALUES (?, ?, ?)",
		id, login, secret,
	)
	if err != nil {
		return fmt.Errorf("ошибка при insert в таблицу: %w", err)
	}
	return nil
}

func (ur *UserRepository) Check(login, pass string) (bool, error) {
	var secret string
	err := ur.db.QueryRow("SELECT secret FROM users WHERE login = ?", login).Scan(&secret)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("ошибка при проверке данных: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(secret), []byte(pass))
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке secret с pass: %w", err)
	}
	return true, nil
}
