package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/domain"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/infrastructure/http/middleware"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/logger"
	"github.com/klimenkokayot/calc-user-go/config"
)

type ProxyServer struct {
	mux     *mux.Router
	handler domain.ProxyHandler
	logger  logger.Logger
	config  *config.Config
}

func (p *ProxyServer) setupRoutes() {
	p.logger.Info("Установка обработчиков.")
	p.mux.PathPrefix("/").Handler(p.handler)
}

func (p *ProxyServer) setupMiddlewares() {
	p.logger.Info("Установка middlewares.")
	p.mux.Use(middleware.RequestMiddleware)
}

func (p *ProxyServer) Run() error {
	p.logger.Info("Запуск прокси-сервера.")
	p.setupMiddlewares()
	p.setupRoutes()
	p.logger.OK("Прокси-сервер успешно запущен.")
	return http.ListenAndServe(fmt.Sprintf(":%d", p.config.ApiGateway.Http.Port), p.mux)
}

func NewProxyServer(handler domain.ProxyHandler, logger logger.Logger, config *config.Config) (domain.ProxyServer, error) {
	logger.Info("Инициализация прокси-сервера.")
	mux := mux.NewRouter()
	logger.OK("Прокси-сервер успешно инициализирован.")
	return &ProxyServer{
		mux:     mux,
		handler: handler,
		logger:  logger,
		config:  config,
	}, nil
}
