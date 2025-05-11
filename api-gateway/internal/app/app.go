package app

import (
	"github.com/klimenkokayot/calc-user-go/api-gateway/config"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/domain"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/logger"
)

type ProxyApplication struct {
	server domain.ProxyServer
	logger logger.Logger
	config *config.Config
}

func (p *ProxyApplication) Run() error {
	return p.server.Run()
}

func NewProxyApplication(server domain.ProxyServer, logger logger.Logger, config *config.Config) (domain.ProxyApplication, error) {
	return &ProxyApplication{
		server: server,
		logger: logger,
		config: config,
	}, nil
}
