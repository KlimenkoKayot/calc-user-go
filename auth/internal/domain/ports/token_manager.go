package ports

import "github.com/golang-jwt/jwt"

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type TokenManager interface {
	NewAccessToken(values map[string]interface{}) (string, error)
	NewRefreshToken(values map[string]interface{}) (string, error)
	NewTokenPair(accessData map[string]interface{}, refreshData map[string]interface{}) (string, string, error)
	ParseWithClaims(tokenString string) (*jwt.MapClaims, error)
	UpdateTokenPair(refreshToken string) (string, string, error)
	ValidateToken(tokenString string) (bool, error)
	ValidateTokenExpiration(token string) (bool, error)
}
