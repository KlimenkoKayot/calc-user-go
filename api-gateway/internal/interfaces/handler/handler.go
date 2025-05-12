package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/klimenkokayot/calc-user-go/api-gateway/internal/domain"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/logger"
	"github.com/klimenkokayot/calc-user-go/config"
)

type ProxyHandler struct {
	service domain.ProxyService
	logger  logger.Logger
	config  *config.Config
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Получаем целевой URL из сервиса
	targetPath, err := p.service.Proxy(r.URL.Path)
	if err != nil {
		p.logger.Warn("Ошибка при обработке запроса.",
			logger.Field{Key: "ip", Value: r.RemoteAddr},
			logger.Field{Key: "request_path", Value: r.URL.Path},
		)
		http.Error(w, "Ошибка при обработке запроса.", http.StatusInternalServerError)
		return
	}

	// Парсим в URL
	targetURL, err := url.Parse(targetPath)
	if err != nil {
		p.logger.Warn("Ошибка при парсинге targetURL.",
			logger.Field{Key: "ip", Value: r.RemoteAddr},
			logger.Field{Key: "target_path", Value: targetPath},
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Создаем прокси
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	p.logger.Debug("Проксирование запроса",
		logger.Field{Key: "method", Value: r.Method},
		logger.Field{Key: "request_path", Value: r.URL.Path},
		logger.Field{Key: "target_url", Value: targetURL.String()},
	)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL = targetURL
		req.Method = r.Method
		req.Body = r.Body
		req.ContentLength = r.ContentLength
		req.Header = r.Header.Clone()
	}

	// Настраиваем обработчики
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		p.logger.Warn("Ошибка во время запроса",
			logger.Field{Key: "ip", Value: r.RemoteAddr},
			logger.Field{Key: "error", Value: err.Error()},
		)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	proxy.ModifyResponse = func(res *http.Response) error {
		p.logger.Info("Proxy response",
			logger.Field{Key: "status", Value: res.Status},
			logger.Field{Key: "status_code", Value: res.StatusCode},
			logger.Field{Key: "content_type", Value: res.Header.Get("Content-Type")},
		)
		return nil
	}

	proxy.ServeHTTP(w, r)
}

func NewProxyHandler(service domain.ProxyService, logger logger.Logger, config *config.Config) (domain.ProxyHandler, error) {
	logger.Info("Инициализация обработчика.")
	logger.OK("Обработчик успешно инициализирован.")
	return &ProxyHandler{
		service: service,
		logger:  logger,
		config:  config,
	}, nil
}
