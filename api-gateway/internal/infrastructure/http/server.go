package server

import (
	"fmt"
	"net/http"

	"github.com/klimenkokayot/calc-user-go/api-gateway/config"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/domain"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/infrastructure/http/middleware"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/logger"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/router"
)

type ProxyServer struct {
	mux     router.Router
	handler domain.ProxyHandler
	logger  logger.Logger
	config  *config.Config
}

func (p *ProxyServer) setupRoutes() {
	p.mux.POST("/api/v1/register", p.handler.Proxy)
	p.mux.GET("/api/v1/register", p.handler.Proxy)
	p.mux.OPTIONS("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
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
	mux, err := router.NewAdapter(&router.Config{
		Name: config.Router,
	})
	if err != nil {
		return nil, err
	}
	return &ProxyServer{
		mux:     mux,
		handler: handler,
		logger:  logger,
		config:  config,
	}, nil
}
