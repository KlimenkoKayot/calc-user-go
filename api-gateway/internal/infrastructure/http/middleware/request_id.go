package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := uuid.NewString()
		r.Header.Set("request_id", uuid)
		next.ServeHTTP(w, r)
	})
}
