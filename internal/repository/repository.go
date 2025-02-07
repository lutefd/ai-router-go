package repository

import (
	"context"

	"github.com/lutefd/ai-router-go/internal/models"
)

type AIRepositoryInterface interface {
	GenerateContentStream(ctx context.Context, model string, prompt string,
		callback func(string)) error
}

type UserRepositoryInterface interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID string) error
	ListUsers(ctx context.Context) ([]*models.User, error)
	GetUsersChatList(ctx context.Context, userID string) ([]*models.UserChat, error)
}

type ChatRepositoryInterface interface {
	CreateChat(ctx context.Context, chat *models.Chat) error
	DeleteChat(ctx context.Context, chatID string) error
	GetChat(ctx context.Context, chatID string) (*models.Chat, error)
	UpdateChat(ctx context.Context, chat *models.Chat) error
}
