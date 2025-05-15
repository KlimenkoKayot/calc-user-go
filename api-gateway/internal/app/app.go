package app

import (
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/domain"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/logger"
	"github.com/klimenkokayot/calc-user-go/config"
)

type ProxyApplication struct {
	server domain.ProxyServer
	logger logger.Logger
	config *config.Config
}

func (p *ProxyApplication) Run() error {
	p.logger.Info("Запуск приложения.")
	return p.server.Run()
}

func NewProxyApplication(server domain.ProxyServer, logger logger.Logger, config *config.Config) (domain.ProxyApplication, error) {
	logger.Info("Инициализация приложения.")
	logger.OK("Приложение успешно инициализовано.")
	return &ProxyApplication{
		server: server,
		logger: logger,
		config: config,
	}, nil
}
