package handler

import (
	"net/http"

	"github.com/klimenkokayot/calc-user-go/api-gateway/config"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/domain"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/logger"
)

type ProxyHandler struct {
	service domain.ProxyService
	logger  logger.Logger
	config  *config.Config
}

func (p *ProxyHandler) Proxy(w http.ResponseWriter, r *http.Request) {
	path, err := p.service.Proxy(r.URL.Path)
	if err != nil {
		http.Error(w, "Ошибка при проксировании запроса.", http.StatusInternalServerError)
		return
	}
	p.logger.Info("Запрос перенаправлен.", logger.Field{
		Key:   "url",
		Value: path,
	})
	http.Redirect(w, r, path, http.StatusFound)
}

func NewProxyService(service domain.ProxyService, logger logger.Logger, config *config.Config) (domain.ProxyHandler, error) {
	return &ProxyHandler{
		service: service,
		logger:  logger,
		config:  config,
	}, nil
}
