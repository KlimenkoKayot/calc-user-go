package middleware

import (
	"context"
	"net/http"
	"time"
)

func TimeoutMiddleware(read time.Duration, write time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var duration time.Duration
			switch r.Method {
			case http.MethodGet:
				duration = read
			default:
				duration = write
			}
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()

			r = r.WithContext(ctx)

			/*
				TODO TODO TODO TODO TODO TODO
			*/

			// done := make(chan interface{})
			// go func() {
			next.ServeHTTP(w, r)
			// 	close(done)
			// }()

			// select {
			// case <-done:
			// 	return
			// case <-ctx.Done():
			// 	r.Response.StatusCode = http.StatusRequestTimeout
			// 	return
			// }
		})
	}
}
