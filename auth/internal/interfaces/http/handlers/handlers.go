package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/klimenkokayot/avito-go/libs/logger"
	"github.com/klimenkokayot/avito-go/services/auth/internal/domain/model"
	domain "github.com/klimenkokayot/avito-go/services/auth/internal/domain/repository"
	"github.com/klimenkokayot/avito-go/services/auth/internal/domain/service"
	"github.com/klimenkokayot/calc-user-go/config"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      logger.Logger
	cfg         *config.Config
}

func NewAuthHandler(service *service.AuthService, cfg *config.Config, logger logger.Logger) (*AuthHandler, error) {
	logger.Info("Инициализация обработчика.")
	logger.OK("Обработчик успешно инициализирован.")
	return &AuthHandler{
		service,
		logger,
		cfg,
	}, nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Warn("Ошибка при чтении тела запроса.", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(fmt.Errorf("ошибка при чтении тела запроса: %w", err))
		return
	}
	defer r.Body.Close()

	user := &domain.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		h.logger.Warn("Ошибка при парсинге тела запроса.", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(fmt.Errorf("ошибка при парсинге тела запроса: %w", err))
		return
	}

	_, _, err = h.authService.Register(user.Login, user.Secret)
	if err != nil {
		h.logger.Warn("Ошибка при регистрации пользователя", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(fmt.Errorf("ошибка при регистрации: %w", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write([]byte(fmt.Errorf("ошибка при создании AuthHandler: %w", err).Error()))
		return
	}
	defer r.Body.Close()

	user := &domain.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		h.logger.Warn("Ошибка при парсинге тела запроса.", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(fmt.Errorf("ошибка при парсинге тела запроса: %w", err))
		return
	}

	accessToken, refreshToken, err := h.authService.Login(user.Login, user.Secret)
	if err != nil {
		h.logger.Warn("Неудачный вход в аккаунт пользователя", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(fmt.Errorf("ошибка при попытке входа: %w", err))
		return
	}

	err = h.updateTokenPair(w, accessToken, refreshToken)
	if err != nil {
		h.logger.Error("500 !!! Неудачный вход в аккаунт пользователя", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(fmt.Errorf("ошибка при попытке входа: %w", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// func (h *AuthHandler) UpdateTokenPair(ctx context.Context, w )

func (h *AuthHandler) ValidateTokenPair(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write([]byte(fmt.Errorf("ошибка при создании AuthHandler: %w", err).Error()))
		return
	}
	defer r.Body.Close()

	tokenPair := &model.TokenPair{}
	err = json.Unmarshal(body, tokenPair)
	if err != nil {
		h.logger.Warn("Ошибка при парсинге тела запроса.", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(fmt.Errorf("ошибка при парсинге тела запроса: %w", err))
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	valid, err := h.authService.ValidateTokenPair(ctx, tokenPair)
	if err != nil {
		h.logger.Warn("Ошибка валидации токенов", logger.Field{
			Key:   "err",
			Value: err.Error(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(fmt.Errorf("ошибка при валидации токенов: %w", err))
		return
	}

	if valid {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(`
		{
			"is_valid": true
		}
		`)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(`
	{
		"is_valid": false,
		"error": "expired token"
	}
	`)
}

func (h *AuthHandler) updateTokenPair(w http.ResponseWriter, accessToken, refreshToken string) error {
	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: accessToken,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})
	return nil
}
