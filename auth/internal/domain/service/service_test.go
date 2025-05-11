package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/klimenkokayot/avito-go/services/auth/internal/domain/model"
	domain "github.com/klimenkokayot/avito-go/services/auth/internal/domain/repository"
	"github.com/klimenkokayot/avito-go/services/auth/internal/domain/service"
	"github.com/klimenkokayot/avito-go/services/auth/mocks/jwt"
	"github.com/klimenkokayot/avito-go/services/auth/mocks/logger"
	"github.com/klimenkokayot/avito-go/services/auth/mocks/repository"
	"github.com/klimenkokayot/calc-user-go/config"
	"go.uber.org/mock/gomock"
)

func TestNewAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockTokenManager := jwt.NewMockTokenManager(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	cfg := &config.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(1)
	mockLogger.EXPECT().OK("Успешно.").Times(1)

	_, err := service.NewAuthService(mockRepo, mockTokenManager, cfg, mockLogger)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockTokenManager := jwt.NewMockTokenManager(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	cfg := &config.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(1)
	mockLogger.EXPECT().OK("Успешно.").Times(1)

	authService, _ := service.NewAuthService(mockRepo, mockTokenManager, cfg, mockLogger)

	login := "testuser"
	password := "password"

	mockRepo.EXPECT().Add(login, gomock.Any()).Return(nil)
	mockTokenManager.EXPECT().NewTokenPair(
		map[string]interface{}{"login": login},
		map[string]interface{}{},
	).Return("access", "refresh", nil)

	access, refresh, err := authService.Register(login, password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if access != "access" || refresh != "refresh" {
		t.Errorf("Expected tokens 'access' and 'refresh', got '%s' and '%s'", access, refresh)
	}
}

func TestAuthService_Register_UserRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockTokenManager := jwt.NewMockTokenManager(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	cfg := &config.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(1)
	mockLogger.EXPECT().OK("Успешно.").Times(1)

	authService, _ := service.NewAuthService(mockRepo, mockTokenManager, cfg, mockLogger)

	expectedErr := errors.New("user already exists")
	mockRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(expectedErr)

	_, _, err := authService.Register("test", "password")
	if err != expectedErr {
		t.Errorf("Expected error '%v', got '%v'", expectedErr, err)
	}
}

func TestAuthService_Register_TokenGenerationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockTokenManager := jwt.NewMockTokenManager(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	cfg := &config.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(1)
	mockLogger.EXPECT().OK("Успешно.").Times(1)

	authService, _ := service.NewAuthService(mockRepo, mockTokenManager, cfg, mockLogger)

	expectedErr := errors.New("token generation failed")
	mockRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
	mockTokenManager.EXPECT().NewTokenPair(gomock.Any(), gomock.Any()).Return("", "", expectedErr)

	_, _, err := authService.Register("test", "password")
	if err != expectedErr {
		t.Errorf("Expected error '%v', got '%v'", expectedErr, err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockTokenManager := jwt.NewMockTokenManager(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	cfg := &config.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(1)
	mockLogger.EXPECT().OK("Успешно.").Times(1)

	authService, _ := service.NewAuthService(mockRepo, mockTokenManager, cfg, mockLogger)

	login := "testuser"
	password := "password"

	mockRepo.EXPECT().Check(login, gomock.Any()).Return(true, nil)
	mockTokenManager.EXPECT().NewTokenPair(
		map[string]interface{}{"login": login},
		map[string]interface{}{},
	).Return("access", "refresh", nil)

	access, refresh, err := authService.Login(login, password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if access != "access" || refresh != "refresh" {
		t.Errorf("Expected tokens 'access' and 'refresh', got '%s' and '%s'", access, refresh)
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockTokenManager := jwt.NewMockTokenManager(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	cfg := &config.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(1)
	mockLogger.EXPECT().OK("Успешно.").Times(1)

	authService, _ := service.NewAuthService(mockRepo, mockTokenManager, cfg, mockLogger)

	testCases := []struct {
		name     string
		valid    bool
		repoErr  error
		expected error
	}{
		{"InvalidPassword", false, nil, domain.ErrBadPassword},
		{"RepoError", false, errors.New("db error"), domain.ErrBadPassword},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().Check(gomock.Any(), gomock.Any()).Return(tc.valid, tc.repoErr)

			_, _, err := authService.Login("test", "password")
			if err != tc.expected {
				t.Errorf("Expected error '%v', got '%v'", tc.expected, err)
			}
		})
	}
}

func TestAuthService_Login_TokenGenerationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockTokenManager := jwt.NewMockTokenManager(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	cfg := &config.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(1)
	mockLogger.EXPECT().OK("Успешно.").Times(1)

	authService, _ := service.NewAuthService(mockRepo, mockTokenManager, cfg, mockLogger)

	expectedErr := errors.New("token generation failed")
	mockRepo.EXPECT().Check(gomock.Any(), gomock.Any()).Return(true, nil)
	mockTokenManager.EXPECT().NewTokenPair(gomock.Any(), gomock.Any()).Return("", "", expectedErr)

	_, _, err := authService.Login("test", "password")
	if err != expectedErr {
		t.Errorf("Expected error '%v', got '%v'", expectedErr, err)
	}
}

func TestAuthService_ValidateTokenPair_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockTokenManager := jwt.NewMockTokenManager(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	cfg := &config.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(1)
	mockLogger.EXPECT().OK("Успешно.").Times(1)

	authService, _ := service.NewAuthService(mockRepo, mockTokenManager, cfg, mockLogger)

	tokenPair := &model.TokenPair{
		AccessToken:  "valid_access",
		RefreshToken: "valid_refresh",
	}

	mockTokenManager.EXPECT().ValidateTokenExpiration("valid_access").Return(true, nil)
	mockTokenManager.EXPECT().ValidateTokenExpiration("valid_refresh").Return(true, nil)

	valid, err := authService.ValidateTokenPair(context.Background(), tokenPair)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !valid {
		t.Error("Expected tokens to be valid")
	}
}

func TestAuthService_ValidateTokenPair_InvalidTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	mockTokenManager := jwt.NewMockTokenManager(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	cfg := &config.Config{}

	mockLogger.EXPECT().Info(gomock.Any()).Times(1)
	mockLogger.EXPECT().OK("Успешно.").Times(1)

	authService, _ := service.NewAuthService(mockRepo, mockTokenManager, cfg, mockLogger)

	testCases := []struct {
		name          string
		accessValid   bool
		accessErr     error
		refreshValid  bool
		refreshErr    error
		expectedValid bool
		expectedErr   error
	}{
		{"InvalidAccess", false, nil, true, nil, false, nil},
		{"InvalidRefresh", true, nil, false, nil, false, nil},
		{"AccessCheckError", false, errors.New("access error"), true, nil, false, errors.New("access error")},
		{"RefreshCheckError", true, nil, false, errors.New("refresh error"), false, errors.New("refresh error")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenPair := &model.TokenPair{
				AccessToken:  "access",
				RefreshToken: "refresh",
			}

			mockTokenManager.EXPECT().ValidateTokenExpiration("access").Return(tc.accessValid, tc.accessErr)
			if tc.accessErr == nil {
				mockTokenManager.EXPECT().ValidateTokenExpiration("refresh").Return(tc.refreshValid, tc.refreshErr)
			}

			valid, err := authService.ValidateTokenPair(context.Background(), tokenPair)
			if valid != tc.expectedValid {
				t.Errorf("Expected valid %v, got %v", tc.expectedValid, valid)
			}
			if (err != nil && tc.expectedErr == nil) ||
				(err == nil && tc.expectedErr != nil) ||
				(err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error()) {
				t.Errorf("Expected error '%v', got '%v'", tc.expectedErr, err)
			}
		})
	}
}
