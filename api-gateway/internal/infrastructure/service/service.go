package service

import (
	"net/url"

	"github.com/klimenkokayot/calc-user-go/api-gateway/config"
	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/domain"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/logger"
)

type ProxyService struct {
	logger logger.Logger
	config *config.Config
}

func (p *ProxyService) Proxy(path string) (string, error) {
	p.logger.Debug("Проксирование запроса.", logger.Field{
		Key:   "path",
		Value: path,
	})
	switch path {
	case "/api/v1/register":
		p.logger.Debug("Редирект в /api/v1/register")
		return url.JoinPath(p.config.AuthURL, "auth", "register")
	case "/api/v1/login":
		p.logger.Debug("Редирект в /api/v1/login")
		return url.JoinPath(p.config.AuthURL, "auth", "login")
	default:
		p.logger.Debug("Редирект в Calc.")
		return url.JoinPath(p.config.CalcURL, path)
	}
}

func NewProxyService(logger logger.Logger, config *config.Config) (domain.ProxyService, error) {
	logger.Info("Инициализация прокси сервиса.")
	logger.OK("Прокси сервис успешно иницилизирован.")
	return &ProxyService{
		logger: logger,
		config: config,
	}, nil
}
