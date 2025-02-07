package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/lutefd/ai-router-go/internal/mocks"
	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/lutefd/ai-router-go/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserService_GetUsersChatList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := service.NewUserService(mockRepo)

	testChats := []*models.UserChat{
		{
			ID:        "chat-1",
			User:      "user-123",
			ChatTitle: "First Chat",
		},
		{
			ID:        "chat-2",
			User:      "user-123",
			ChatTitle: "Second Chat",
		},
	}

	tests := []struct {
		name    string
		userID  string
		setup   func()
		want    []*models.UserChat
		wantErr bool
	}{
		{
			name:   "successful retrieval",
			userID: "user-123",
			setup: func() {
				mockRepo.EXPECT().
					GetUsersChatList(gomock.Any(), "user-123").
					Return(testChats, nil)
			},
			want:    testChats,
			wantErr: false,
		},
		{
			name:    "empty user ID",
			userID:  "",
			setup:   func() {},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "repository error",
			userID: "user-123",
			setup: func() {
				mockRepo.EXPECT().
					GetUsersChatList(gomock.Any(), "user-123").
					Return(nil, fmt.Errorf("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := userService.GetUsersChatList(context.Background(), tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
