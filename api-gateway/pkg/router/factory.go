package router

import (
	"fmt"

	"github.com/klimenkokayot/calc-user-go/api-gateway/router/adapters/gorilla"
	"github.com/klimenkokayot/calc-user-go/api-gateway/router/domain"
)

const (
	AdapterGorilla = "gorilla"
)

var (
	ErrUnknownAdapter = fmt.Errorf("роутер не поддерижвается")
)

type (
	Router = domain.Router
)

type Config struct {
	Name string
}

func NewAdapter(cfg *Config) (Router, error) {
	switch cfg.Name {
	case AdapterGorilla:
		return gorilla.NewAdapter()
	default:
		return nil, ErrUnknownAdapter
	}
}
