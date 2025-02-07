package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/lutefd/ai-router-go/internal/repository"
)

type AuthService struct {
	userRepo           repository.UserRepositoryInterface
	jwtSecret          []byte
	tokenExpiry        time.Duration
	refreshTokenExpiry time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func NewAuthService(userRepo repository.UserRepositoryInterface, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		jwtSecret:          []byte(jwtSecret),
		tokenExpiry:        15 * time.Minute,
		refreshTokenExpiry: 30 * 24 * time.Hour,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) AuthenticateUser(ctx context.Context, email string, name string, googleID string) (*models.User, *TokenPair, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() != "user not found" {
			return nil, nil, fmt.Errorf("failed to get user: %w", err)
		}

		user = &models.User{
			ID:    googleID,
			Name:  name,
			Email: email,
			Role:  "user",
		}
		if err := s.userRepo.CreateUser(ctx, user); err != nil {
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	tokenPair, err := s.GenerateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokenPair, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user cannot be nil")
	}

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *AuthService) GenerateTokenPair(user *models.User) (*TokenPair, error) {
	accessToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	refreshClaims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: signedRefreshToken,
		ExpiresIn:    int64(s.tokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := s.userRepo.GetUser(context.Background(), claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.GenerateTokenPair(user)
}
