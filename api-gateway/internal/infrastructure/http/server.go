package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/klimenkokayot/calc-user-go/api-gateway/config"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/domain"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/infrastructure/http/middleware"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/logger"
)

type ProxyServer struct {
	mux     *mux.Router
	handler domain.ProxyHandler
	logger  logger.Logger
	config  *config.Config
}

func (p *ProxyServer) setupRoutes() {
	p.mux.PathPrefix("/").Handler(p.handler)
}

func (p *ProxyServer) setupMiddlewares() {
	p.mux.Use(middleware.RequestMiddleware)
}

func (p *ProxyServer) Run() error {
	p.setupMiddlewares()
	p.setupRoutes()
	return http.ListenAndServe(fmt.Sprintf(":%d", p.config.Port), p.mux)
}

func NewProxyServer(handler domain.ProxyHandler, logger logger.Logger, config *config.Config) (domain.ProxyServer, error) {
	mux := mux.NewRouter()
	return &ProxyServer{
		mux:     mux,
		handler: handler,
		logger:  logger,
		config:  config,
	}, nil
}
