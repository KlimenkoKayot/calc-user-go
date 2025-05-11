package middleware

import (
	"net/http"
	"time"

	"github.com/klimenkokayot/avito-go/libs/logger"
	"github.com/klimenkokayot/avito-go/libs/logger/domain"
)

func LoggerMiddleware(logger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			defer func() {
				if err := recover(); err != nil {
					logger.WithFields(
						domain.Field{Key: "method", Value: r.Method},
						domain.Field{Key: "path", Value: r.URL.Path},
						domain.Field{Key: "duration", Value: time.Since(start)},
						domain.Field{Key: "ip", Value: r.RemoteAddr},
					).Error("Перехвачена ошибка в запросе.")
				}
			}()

			next.ServeHTTP(w, r)

			logger.WithFields(
				domain.Field{Key: "method", Value: r.Method},
				domain.Field{Key: "path", Value: r.URL.Path},
				domain.Field{Key: "duration", Value: time.Since(start)},
				domain.Field{Key: "ip", Value: r.RemoteAddr},
			).Info("Обработка запроса.")
		})
	}
}
