package orchestartor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	models "github.com/klimenkokayot/calc-net-go/internal/shared/models"
	logger "github.com/klimenkokayot/calc-net-go/pkg/logger"
)

type AuthClient struct {
	client      http.Client
	logger      logger.Logger
	authBaseURL string
}

func (c *AuthClient) NewAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessTokenCookie, err := r.Cookie("access_token")
			if err != nil {
				http.Error(w, "access_token cookie не найден", http.StatusUnauthorized)
				return
			}
			accessTokenString := accessTokenCookie.String()

			refreshTokenCookie, err := r.Cookie("refresh_token")
			if err != nil {
				http.Error(w, "refresh_token cookie не найден", http.StatusUnauthorized)
				return
			}
			refreshTokenString := refreshTokenCookie.String()

			_, err = c.Authenticate(nil, &models.TokenPair{
				AccessToken:  accessTokenString,
				RefreshToken: refreshTokenString,
			})
			if err != nil {
				http.Error(w, fmt.Sprintf("Не удалось обновить пару токенов, %s", err.Error()), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (a *AuthClient) Authenticate(ctx context.Context, tokenPair *models.TokenPair) (userID string, err error) {
	a.logger.Info("Проверка токенов.")
	data, err := json.Marshal(tokenPair)
	if err != nil {
		a.logger.Warn("Не удалось сформировать тело запроса.")
		return "", err
	}
	buffer := new(bytes.Buffer)
	buffer.WriteString(string(data))
	authValidatePath, err := url.JoinPath(a.authBaseURL, "auth", "validate")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, authValidatePath, buffer)
	if err != nil {
		a.logger.Warn("Не удалось сформировать http.Request.")
		return "", err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		a.logger.Warn("Не удалось отравить запрос.")
		return "", err
	}
	defer resp.Body.Close()

	// Проверка на верификацию
	if resp.StatusCode == http.StatusOK {
		a.logger.OK("Токен валиден.")
		return "todo_user_id", nil
	}
	a.logger.OK("Токен не валиден.")
	return "", nil
}

func NewAuthClient(authBaseURL string, logger logger.Logger) *AuthClient {
	return &AuthClient{
		client:      http.Client{},
		logger:      logger,
		authBaseURL: authBaseURL,
	}
}
