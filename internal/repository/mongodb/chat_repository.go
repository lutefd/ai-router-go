package mongodb

import (
	"context"
	"fmt"

	"github.com/lutefd/ai-router-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepository struct {
	db *mongo.Database
}

func NewChatRepository(db *mongo.Database) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateChat(ctx context.Context, chat *models.Chat) error {
	_, err := r.db.Collection("chats").InsertOne(ctx, chat)
	if err != nil {
		return fmt.Errorf("failed to create chat: %w", err)
	}
	return nil
}

func (r *ChatRepository) DeleteChat(ctx context.Context, chatID string) error {
	result, err := r.db.Collection("chats").DeleteOne(ctx, bson.M{"_id": chatID})
	if err != nil {
		return fmt.Errorf("error deleting chat: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("chat not found")
	}
	return nil
}

func (r *ChatRepository) GetChat(ctx context.Context, id string) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.Collection("chats").FindOne(ctx, bson.M{"_id": id}).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("chat not found")
		}
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}
	return &chat, nil
}

func (r *ChatRepository) UpdateChat(ctx context.Context, chat *models.Chat) error {
	result, err := r.db.Collection("chats").ReplaceOne(ctx, bson.M{"_id": chat.ID}, chat)
	if err != nil {
		return fmt.Errorf("error updating chat: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("chat not found")
	}
	return nil
}
