package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/klimenkokayot/avito-go/libs/logger"
	"github.com/klimenkokayot/avito-go/libs/router"
	"github.com/klimenkokayot/avito-go/services/auth/config"
	"github.com/klimenkokayot/avito-go/services/auth/internal/domain"
	"github.com/klimenkokayot/avito-go/services/auth/internal/infrastructure/http/middleware"
	"github.com/klimenkokayot/avito-go/services/auth/internal/interfaces/http/handlers"
)

type AuthServer struct {
	handler      *handlers.AuthHandler
	router       router.Router
	logger       logger.Logger
	readTimeout  time.Duration
	writeTimeout time.Duration
	cfg          *config.Config
}

func NewAuthServer(handler *handlers.AuthHandler, cfg *config.Config, logger logger.Logger) (domain.Server, error) {
	logger.Info("Инициализация сервера.")
	router, err := router.NewAdapter(&router.Config{
		Name: cfg.Router,
	})
	if err != nil {
		return nil, err
	}

	server := &AuthServer{
		handler,
		router,
		logger,
		cfg.Server.ReadTimeout,
		cfg.Server.WriteTimeout,
		cfg,
	}

	err = server.setupRoutes()
	if err != nil {
		return nil, err
	}

	err = server.setupMiddleware()
	if err != nil {
		return nil, err
	}

	logger.OK("Успешно.")
	return server, nil
}

func (a *AuthServer) setupRoutes() error {
	a.logger.Info("Инициализация ручек.")

	// Аутентификация пользователя
	a.router.POST("/auth/login", a.handler.Login)
	a.router.OPTIONS("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Регистрация пользователя
	a.router.POST("/auth/register", a.handler.Register)
	a.router.OPTIONS("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Валидация токенов
	a.router.GET("/auth/validate", a.handler.ValidateTokenPair)
	a.router.OPTIONS("/auth/validate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	a.logger.OK("Успешно.")
	return nil
}

func (a *AuthServer) setupMiddleware() error {
	a.router.Use(middleware.CorsMiddleware())
	a.router.Use(middleware.LoggerMiddleware(a.logger))
	a.router.Use(middleware.TimeoutMiddleware(a.readTimeout, a.writeTimeout))
	return nil
}

func (a *AuthServer) Run() error {
	a.logger.Info("Сервер запущен.", logger.Field{Key: "port", Value: a.cfg.Server.Port})
	return http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.Server.Port), a.router)
}
