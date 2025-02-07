package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/lutefd/ai-router-go/internal/repository/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestChatRepository_CreateChat(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewChatRepository(conn.DB)
	ctx := context.Background()

	tests := []struct {
		name    string
		chat    *models.Chat
		wantErr bool
	}{
		{
			name: "successful chat creation",
			chat: &models.Chat{
				ID:        "chat-id",
				Title:     "Test Chat",
				Messages:  []models.Message{{ID: "msg-1", Text: "Hello", Role: "user", SentAt: time.Now()}},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateChat(ctx, tt.chat)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			var found models.Chat
			err = conn.DB.Collection("chats").FindOne(ctx, bson.M{"_id": tt.chat.ID}).Decode(&found)
			require.NoError(t, err)
			assert.Equal(t, tt.chat.ID, found.ID)
			assert.Equal(t, tt.chat.Title, found.Title)
		})
	}
}

func TestChatRepository_GetChat(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewChatRepository(conn.DB)
	ctx := context.Background()

	testChat := &models.Chat{
		ID:        "chat-1",
		Title:     "Test Chat",
		Messages:  []models.Message{{ID: "msg-1", Text: "Hello", Role: "user", SentAt: time.Now()}},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.CreateChat(ctx, testChat)
	require.NoError(t, err)

	tests := []struct {
		name    string
		chatID  string
		want    *models.Chat
		wantErr bool
	}{
		{
			name:    "existing chat",
			chatID:  "chat-1",
			want:    testChat,
			wantErr: false,
		},
		{
			name:    "non-existent chat",
			chatID:  "invalid-id",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chat, err := repo.GetChat(ctx, tt.chatID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.ID, chat.ID)
			assert.Equal(t, tt.want.Title, chat.Title)
		})
	}
}
