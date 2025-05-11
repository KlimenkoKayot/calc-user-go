package app

import (
	"fmt"

	"github.com/klimenkokayot/avito-go/libs/jwt"
	"github.com/klimenkokayot/avito-go/libs/logger"
	"github.com/klimenkokayot/avito-go/services/auth/config"
	"github.com/klimenkokayot/avito-go/services/auth/internal/domain"
	"github.com/klimenkokayot/avito-go/services/auth/internal/domain/service"
	server "github.com/klimenkokayot/avito-go/services/auth/internal/infrastructure/http"
	repo "github.com/klimenkokayot/avito-go/services/auth/internal/infrastructure/http/repository"
	"github.com/klimenkokayot/avito-go/services/auth/internal/interfaces/http/handlers"
)

type Application struct {
	server domain.Server
	logger logger.Logger
	config *config.Config
}

func NewApplication(cfg *config.Config, logger logger.Logger) (domain.Application, error) {
	logger.Info("Инициализация application.")

	repoLogger := logger.WithLayer("REPO")
	userRepo, err := repo.NewUserRepository(cfg, repoLogger)
	if err != nil {
		return nil, wrapError("инициализации репозитория", err)
	}

	tokenManager, err := jwt.NewTokenManager(cfg.JwtSecretKey, cfg.AccessTokenExpiration, cfg.RefreshTokenExpiration)
	if err != nil {
		return nil, err
	}

	serviceLogger := logger.WithLayer("SERVICE")
	authService, err := service.NewAuthService(userRepo, tokenManager, cfg, serviceLogger)
	if err != nil {
		return nil, wrapError("инициализации сервиса:", err)
	}

	handlerLogger := logger.WithLayer("HANDLER")
	handler, err := handlers.NewAuthHandler(authService, cfg, handlerLogger)
	if err != nil {
		return nil, wrapError("инициализации обработчиков:", err)
	}

	serverLogger := logger.WithLayer("SERVER")
	server, err := server.NewAuthServer(handler, cfg, serverLogger)
	if err != nil {
		return nil, wrapError("инициализации сервера:", err)
	}

	logger.OK("Application успешно инициализирован.")
	return &Application{
		server,
		logger,
		cfg,
	}, nil
}

func (a *Application) Run() error {
	a.logger.Info("Запуск application.")
	return a.server.Run()
}

func wrapError(msg string, err error) error {
	return fmt.Errorf("Ошибка %s: %w", msg, err)
}
