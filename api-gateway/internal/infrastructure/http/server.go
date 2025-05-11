package server

import (
	"github.com/klimenkokayot/calc-user-go/api-gateway/config"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/domain"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/logger"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/router"
)

type ProxyServer struct {
	mux     router.Router
	handler domain.ProxyHandler
	logger  logger.Logger
	config  *config.Config
}

func (p *ProxyServer) Run() error {

}

func NewProxyService(handler domain.ProxyHandler, logger logger.Logger, config *config.Config) domain.ProxyServer {
	return &ProxyServer{
		handler: handler,
		logger:  logger,
		config:  config,
	}
}
