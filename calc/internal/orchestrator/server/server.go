package orchestrator

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	client "github.com/klimenkokayot/calc-net-go/internal/orchestrator/client"
	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	handler "github.com/klimenkokayot/calc-net-go/internal/orchestrator/server/handler"
)

// Структура сервера
type Server struct {
	config  *config.Config
	handler *handler.OrchestratorHandler
	client  *client.AuthClient
	mux     *mux.Router
}

// Создание нового экземпляра сервера
func NewServer(config *config.Config, handler *handler.OrchestratorHandler, client *client.AuthClient) (*Server, error) {
	mux := mux.NewRouter()
	return &Server{
		config:  config,
		handler: handler,
		client:  client,
		mux:     mux,
	}, nil
}

// Запуск сервера, использует роутер gorilla/mux
func (s *Server) Run() error {
	err := s.setupMiddlewares()
	if err != nil {
		return err
	}

	err = s.setupRoutes()
	if err != nil {
		return err
	}

	log.Printf("Server started at port :%d\n", s.config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s.mux); err != nil {
		return err
	}
	return nil
}

func (s *Server) setupMiddlewares() error {
	authMiddleware := s.client.NewAuthMiddleware()
	s.mux.Use(authMiddleware)
	return nil
}

func (s *Server) setupRoutes() error {
	// Разные endpoint`ы
	s.mux.HandleFunc("/", s.handler.Index)
	s.mux.HandleFunc("/api/v1/calculate", s.handler.NewExpression)
	s.mux.HandleFunc("/api/v1/expressions", s.handler.Expressions)
	s.mux.HandleFunc("/api/v1/expressions/{id}", s.handler.Expression)
	s.mux.HandleFunc("/internal/task", s.handler.PostTask).Methods("POST")
	s.mux.HandleFunc("/internal/task", s.handler.GetTask).Methods("GET")

	staticDir := filepath.Join(".", "web", "static")
	fs := http.FileServer(http.Dir(staticDir))
	s.mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	return nil
}
