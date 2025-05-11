package gorilla

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/klimenkokayot/calc-user-go/api-gateway/pkg/router/domain"
)

type AdapterGorilla struct {
	router *mux.Router
}

func NewAdapter() (domain.Router, error) {
	mux := mux.NewRouter()
	return &AdapterGorilla{
		mux,
	}, nil
}

func (a *AdapterGorilla) GET(path string, handler domain.HandlerFunc) {
	a.router.HandleFunc(path, handler).Methods(http.MethodGet)
}

func (a *AdapterGorilla) POST(path string, handler domain.HandlerFunc) {
	a.router.HandleFunc(path, handler).Methods(http.MethodPost)
}

func (a *AdapterGorilla) OPTIONS(path string, handler domain.HandlerFunc) {
	a.router.HandleFunc(path, handler).Methods(http.MethodOptions)
}

func (a *AdapterGorilla) Handle(path string, handler domain.Handler) {
	a.router.Handle(path, handler)
}

func (a *AdapterGorilla) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *AdapterGorilla) Use(middleware domain.MiddlewareFunc) {
	a.router.Use(middleware)
}
