package orchestrator

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	client "github.com/klimenkokayot/calc-net-go/internal/orchestrator/client"
	handler "github.com/klimenkokayot/calc-net-go/internal/orchestrator/server/handler"
	"github.com/klimenkokayot/calc-user-go/config"
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
	authRouter := s.mux.PathPrefix("/").Subrouter()
	authRouter.Use(s.client.NewAuthMiddleware())

	s.setupAgentRoutes()

	s.setupAuthRoutes(authRouter)

	s.setupStaticRoute()

	log.Printf("Server started at port :%d\n", s.config.Calc.Orchestrator.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.config.Calc.Orchestrator.Port), s.mux); err != nil {
		return err
	}
	return nil
}

// Роуты для внутреннего взаимодействия
func (s *Server) setupAgentRoutes() {
	agentRouter := s.mux.PathPrefix("/internal").Subrouter()
	agentRouter.HandleFunc("/task", s.handler.PostTask).Methods("POST")
	agentRouter.HandleFunc("/task", s.handler.GetTask).Methods("GET")
}

// Основные роуты с аутентификацией
func (s *Server) setupAuthRoutes(authRouter *mux.Router) {
	authRouter.HandleFunc("/", s.handler.Index)
	authRouter.HandleFunc("/api/v1/calculate", s.handler.NewExpression).Methods("POST")
	authRouter.HandleFunc("/api/v1/expressions", s.handler.Expressions).Methods("GET")
	authRouter.HandleFunc("/api/v1/expressions/{id}", s.handler.Expression).Methods("GET")
}
func (s *Server) setupStaticRoute() {
	staticDir := filepath.Join(".", "web", "static")
	fs := http.FileServer(http.Dir(staticDir))
	s.mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
}
