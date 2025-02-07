package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lutefd/ai-router-go/internal/mocks"
	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/lutefd/ai-router-go/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthService_AuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	authService := service.NewAuthService(mockUserRepo, "test-secret")

	tests := []struct {
		name      string
		email     string
		userName  string
		googleID  string
		setupMock func()
		wantUser  *models.User
		wantErr   bool
	}{
		{
			name:     "existing user successful authentication",
			email:    "test@example.com",
			userName: "Test User",
			googleID: "123",
			setupMock: func() {
				mockUserRepo.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Return(&models.User{
						ID:    "123",
						Email: "test@example.com",
						Name:  "Test User",
						Role:  "user",
					}, nil)
			},
			wantUser: &models.User{
				ID:    "123",
				Email: "test@example.com",
				Name:  "Test User",
				Role:  "user",
			},
			wantErr: false,
		},
		{
			name:     "new user successful creation",
			email:    "new@example.com",
			userName: "New User",
			googleID: "456",
			setupMock: func() {
				mockUserRepo.EXPECT().
					GetUserByEmail(gomock.Any(), "new@example.com").
					Return(nil, fmt.Errorf("user not found"))
				mockUserRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantUser: &models.User{
				ID:    "456",
				Email: "new@example.com",
				Name:  "New User",
				Role:  "user",
			},
			wantErr: false,
		},
		{
			name:     "repository error",
			email:    "error@example.com",
			userName: "Error User",
			googleID: "789",
			setupMock: func() {
				mockUserRepo.EXPECT().
					GetUserByEmail(gomock.Any(), "error@example.com").
					Return(nil, fmt.Errorf("database error"))
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			user, tokenPair, err := authService.AuthenticateUser(context.Background(), tt.email, tt.userName, tt.googleID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantUser, user)
			assert.NotNil(t, tokenPair)
			assert.NotEmpty(t, tokenPair.AccessToken)
			assert.NotEmpty(t, tokenPair.RefreshToken)
			assert.Equal(t, int64(15*60), tokenPair.ExpiresIn)
			claims, err := authService.ValidateToken(tokenPair.AccessToken)
			require.NoError(t, err)
			assert.Equal(t, user.ID, claims.UserID)
			assert.Equal(t, user.Email, claims.Email)
			assert.Equal(t, user.Role, claims.Role)

			refreshClaims, err := authService.ValidateToken(tokenPair.RefreshToken)
			require.NoError(t, err)
			assert.Equal(t, user.ID, refreshClaims.UserID)
			assert.Equal(t, user.Email, refreshClaims.Email)
		})
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	authService := service.NewAuthService(mockUserRepo, "test-secret")

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &models.User{
				ID:    "123",
				Email: "test@example.com",
				Name:  "Test User",
				Role:  "user",
			},
			wantErr: false,
		},
		{
			name:    "nil user",
			user:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := authService.GenerateToken(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, token)

			claims, err := authService.ValidateToken(token)
			require.NoError(t, err)
			assert.Equal(t, tt.user.ID, claims.UserID)
			assert.Equal(t, tt.user.Email, claims.Email)
			assert.Equal(t, tt.user.Name, claims.Name)
			assert.Equal(t, tt.user.Role, claims.Role)
			assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
		})
	}
}

func TestAuthService_ValidateToken_Invalid(t *testing.T) {
	authService := service.NewAuthService(nil, "test-secret")

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "malformed token",
			token: "invalid.token.string",
		},
		{
			name:  "wrong signature",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.wrong-signature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := authService.ValidateToken(tt.token)
			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}

func TestAuthService_GenerateTokenPair(t *testing.T) {
	authService := service.NewAuthService(nil, "test-secret")

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &models.User{
				ID:    "123",
				Email: "test@example.com",
				Name:  "Test User",
				Role:  "user",
			},
			wantErr: false,
		},
		{
			name:    "nil user",
			user:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenPair, err := authService.GenerateTokenPair(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, tokenPair.AccessToken)
			assert.NotEmpty(t, tokenPair.RefreshToken)
			assert.Equal(t, int64(15*60), tokenPair.ExpiresIn)

			claims, err := authService.ValidateToken(tokenPair.AccessToken)
			require.NoError(t, err)
			assert.Equal(t, tt.user.ID, claims.UserID)
			assert.Equal(t, tt.user.Email, claims.Email)
			assert.Equal(t, tt.user.Role, claims.Role)
			assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
			assert.True(t, claims.ExpiresAt.Time.Before(time.Now().Add(16*time.Minute)))

			refreshClaims, err := authService.ValidateToken(tokenPair.RefreshToken)
			require.NoError(t, err)
			assert.Equal(t, tt.user.ID, refreshClaims.UserID)
			assert.Equal(t, tt.user.Email, refreshClaims.Email)
			assert.True(t, refreshClaims.ExpiresAt.Time.After(time.Now().Add(29*24*time.Hour)))
		})
	}
}

func TestAuthService_RefreshAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	authService := service.NewAuthService(mockUserRepo, "test-secret")

	user := &models.User{
		ID:    "123",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "user",
	}
	initialTokenPair, err := authService.GenerateTokenPair(user)
	require.NoError(t, err)

	tests := []struct {
		name         string
		refreshToken string
		setupMock    func()
		wantErr      bool
	}{
		{
			name:         "valid refresh token",
			refreshToken: initialTokenPair.RefreshToken,
			setupMock: func() {
				mockUserRepo.EXPECT().
					GetUser(gomock.Any(), "123").
					Return(user, nil)
			},
			wantErr: false,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid-token",
			setupMock:    func() {},
			wantErr:      true,
		},
		{
			name:         "user not found",
			refreshToken: initialTokenPair.RefreshToken,
			setupMock: func() {
				mockUserRepo.EXPECT().
					GetUser(gomock.Any(), "123").
					Return(nil, fmt.Errorf("user not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			newTokenPair, err := authService.RefreshAccessToken(tt.refreshToken)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, newTokenPair.AccessToken)
			assert.NotEmpty(t, newTokenPair.RefreshToken)
			assert.Equal(t, int64(15*60), newTokenPair.ExpiresIn)
			assert.NotEqual(t, initialTokenPair.AccessToken, newTokenPair.AccessToken)
			assert.NotEqual(t, initialTokenPair.RefreshToken, newTokenPair.RefreshToken)

			claims, err := authService.ValidateToken(newTokenPair.AccessToken)
			require.NoError(t, err)
			assert.Equal(t, user.ID, claims.UserID)
			assert.Equal(t, user.Email, claims.Email)
			assert.Equal(t, user.Role, claims.Role)
		})
	}
}
