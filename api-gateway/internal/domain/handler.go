package domain

import "net/http"

type ProxyHandler interface {
	Proxy(w http.ResponseWriter, r *http.Request)
}
