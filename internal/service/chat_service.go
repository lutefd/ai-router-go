package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/lutefd/ai-router-go/internal/repository"
	"github.com/lutefd/ai-router-go/pkg/idgen"
)

type ChatService struct {
	chatRepo repository.ChatRepositoryInterface
}

func NewChatService(chatRepo repository.ChatRepositoryInterface) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

func (s *ChatService) CreateChat(ctx context.Context, chat *models.Chat) error {
	if chat.Title == "" {
		return fmt.Errorf("chat title is required")
	}

	chat.ID = generateID()
	chat.CreatedAt = time.Now()
	chat.UpdatedAt = time.Now()
	chat.Messages = []models.Message{}

	return s.chatRepo.CreateChat(ctx, chat)
}

func (s *ChatService) GetChat(ctx context.Context, id string) (*models.Chat, error) {
	if id == "" {
		return nil, fmt.Errorf("chat ID is required")
	}

	return s.chatRepo.GetChat(ctx, id)
}

func (s *ChatService) UpdateChat(ctx context.Context, chat *models.Chat) error {
	if chat.ID == "" {
		return fmt.Errorf("chat ID is required")
	}
	if chat.Title == "" {
		return fmt.Errorf("chat title is required")
	}

	existingChat, err := s.chatRepo.GetChat(ctx, chat.ID)
	if err != nil {
		return fmt.Errorf("failed to get chat: %w", err)
	}

	existingChat.Title = chat.Title
	existingChat.UpdatedAt = time.Now()

	return s.chatRepo.UpdateChat(ctx, existingChat)
}

func (s *ChatService) DeleteChat(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("chat ID is required")
	}

	return s.chatRepo.DeleteChat(ctx, id)
}

func generateID() string {
	return idgen.Generate()
}
