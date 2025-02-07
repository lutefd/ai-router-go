package service_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/lutefd/ai-router-go/internal/mocks"
	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/lutefd/ai-router-go/internal/service"
	"github.com/lutefd/ai-router-go/pkg/idgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	if err := idgen.Init(1); err != nil {
		log.Fatal("failed to initialize snowflake node:", err)
	}
	os.Exit(m.Run())
}

func TestChatService_CreateChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepositoryInterface(ctrl)
	chatService := service.NewChatService(mockRepo)

	tests := []struct {
		name    string
		chat    *models.Chat
		setup   func()
		wantErr bool
	}{
		{
			name: "successful creation",
			chat: &models.Chat{
				Title: "Test Chat",
				User:  "user-123",
			},
			setup: func() {
				mockRepo.EXPECT().
					CreateChat(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "empty title",
			chat: &models.Chat{
				Title: "",
				User:  "user-123",
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "repository error",
			chat: &models.Chat{
				Title: "Test Chat",
				User:  "user-123",
			},
			setup: func() {
				mockRepo.EXPECT().
					CreateChat(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := chatService.CreateChat(context.Background(), tt.chat)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, tt.chat.ID)
			assert.NotEmpty(t, tt.chat.CreatedAt)
			assert.NotEmpty(t, tt.chat.UpdatedAt)
			assert.Empty(t, tt.chat.Messages)
		})
	}
}

func TestChatService_GetChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepositoryInterface(ctrl)
	chatService := service.NewChatService(mockRepo)

	testChat := &models.Chat{
		ID:        "chat-123",
		Title:     "Test Chat",
		User:      "user-123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		chatID  string
		setup   func()
		want    *models.Chat
		wantErr bool
	}{
		{
			name:   "successful retrieval",
			chatID: "chat-123",
			setup: func() {
				mockRepo.EXPECT().
					GetChat(gomock.Any(), "chat-123").
					Return(testChat, nil)
			},
			want:    testChat,
			wantErr: false,
		},
		{
			name:    "empty ID",
			chatID:  "",
			setup:   func() {},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "not found",
			chatID: "nonexistent",
			setup: func() {
				mockRepo.EXPECT().
					GetChat(gomock.Any(), "nonexistent").
					Return(nil, fmt.Errorf("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := chatService.GetChat(context.Background(), tt.chatID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChatService_UpdateChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepositoryInterface(ctrl)
	chatService := service.NewChatService(mockRepo)

	existingChat := &models.Chat{
		ID:        "chat-123",
		Title:     "Original Title",
		User:      "user-123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		chat    *models.Chat
		setup   func()
		wantErr bool
	}{
		{
			name: "successful update",
			chat: &models.Chat{
				ID:    "chat-123",
				Title: "Updated Title",
			},
			setup: func() {
				mockRepo.EXPECT().
					GetChat(gomock.Any(), "chat-123").
					Return(existingChat, nil)
				mockRepo.EXPECT().
					UpdateChat(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			chat: &models.Chat{
				Title: "Updated Title",
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "empty title",
			chat: &models.Chat{
				ID: "chat-123",
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "not found",
			chat: &models.Chat{
				ID:    "nonexistent",
				Title: "Updated Title",
			},
			setup: func() {
				mockRepo.EXPECT().
					GetChat(gomock.Any(), "nonexistent").
					Return(nil, fmt.Errorf("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := chatService.UpdateChat(context.Background(), tt.chat)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestChatService_DeleteChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepositoryInterface(ctrl)
	chatService := service.NewChatService(mockRepo)

	tests := []struct {
		name    string
		chatID  string
		setup   func()
		wantErr bool
	}{
		{
			name:   "successful deletion",
			chatID: "chat-123",
			setup: func() {
				mockRepo.EXPECT().
					DeleteChat(gomock.Any(), "chat-123").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "empty ID",
			chatID:  "",
			setup:   func() {},
			wantErr: true,
		},
		{
			name:   "not found",
			chatID: "nonexistent",
			setup: func() {
				mockRepo.EXPECT().
					DeleteChat(gomock.Any(), "nonexistent").
					Return(fmt.Errorf("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := chatService.DeleteChat(context.Background(), tt.chatID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
