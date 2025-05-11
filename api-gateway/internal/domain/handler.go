package domain

import "net/http"

type ProxyHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
