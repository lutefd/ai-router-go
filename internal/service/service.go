package service

import (
	"context"

	"github.com/lutefd/ai-router-go/internal/models"
)

type AIServiceInterface interface {
	GenerateResponse(ctx context.Context, model string, prompt string,
		callback func(string)) error
	GenerateOpenAIResponse(ctx context.Context, model string, prompt string,
		callback func(string)) error
	GenerateDeepSeekResponse(ctx context.Context, model string, prompt string,
		callback func(string)) error
}

type AuthServiceInterface interface {
	AuthenticateUser(ctx context.Context, email string, name string, googleID string) (*models.User, string, error)
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}
