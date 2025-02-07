package service

import (
	"context"
	"fmt"

	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/lutefd/ai-router-go/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserService(userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUsersChatList(ctx context.Context, userID string) ([]*models.UserChat, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	return s.userRepo.GetUsersChatList(ctx, userID)
}
