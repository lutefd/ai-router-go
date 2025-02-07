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
	userRepo    repository.UserRepositoryInterface
	jwtSecret   []byte
	tokenExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepositoryInterface, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: 24 * time.Hour,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) AuthenticateUser(ctx context.Context, email string, name string, googleID string) (*models.User, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() != "user not found" {
			return nil, "", fmt.Errorf("failed to get user: %w", err)
		}

		user = &models.User{
			ID:    googleID,
			Name:  name,
			Email: email,
			Role:  "user",
		}
		if err := s.userRepo.CreateUser(ctx, user); err != nil {
			return nil, "", fmt.Errorf("failed to create user: %w", err)
		}
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user cannot be nil")
	}

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
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
