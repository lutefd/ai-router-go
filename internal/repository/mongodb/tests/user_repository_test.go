package mongodb_test

import (
	"context"
	"testing"

	"github.com/lutefd/ai-router-go/internal/models"
	"github.com/lutefd/ai-router-go/internal/repository/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUserRepository_CreateUser(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewUserRepository(conn.DB)
	ctx := context.Background()

	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "successful user creation",
			user: &models.User{
				ID:    "test-id",
				Name:  "Test User",
				Email: "test@example.com",
				Role:  "user",
			},
			wantErr: false,
		},
		{
			name: "duplicate user",
			user: &models.User{
				ID:    "test-id",
				Name:  "Test User",
				Email: "test@example.com",
				Role:  "user",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateUser(ctx, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			var found models.User
			err = conn.DB.Collection("users").FindOne(ctx, bson.M{"_id": tt.user.ID}).Decode(&found)
			require.NoError(t, err)
			assert.Equal(t, tt.user.ID, found.ID)
			assert.Equal(t, tt.user.Email, found.Email)
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewUserRepository(conn.DB)
	ctx := context.Background()

	testUser := &models.User{
		ID:    "test-id",
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "user",
	}
	err := repo.CreateUser(ctx, testUser)
	require.NoError(t, err)

	tests := []struct {
		name     string
		userID   string
		wantUser *models.User
		wantErr  bool
	}{
		{
			name:     "existing user",
			userID:   "test-id",
			wantUser: testUser,
			wantErr:  false,
		},
		{
			name:     "non-existent user",
			userID:   "invalid-id",
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetUser(ctx, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantUser.ID, user.ID)
			assert.Equal(t, tt.wantUser.Email, user.Email)
		})
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewUserRepository(conn.DB)
	ctx := context.Background()

	testUser := &models.User{
		ID:    "test-id",
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "user",
	}
	err := repo.CreateUser(ctx, testUser)
	require.NoError(t, err)

	tests := []struct {
		name     string
		email    string
		wantUser *models.User
		wantErr  bool
	}{
		{
			name:     "existing email",
			email:    "test@example.com",
			wantUser: testUser,
			wantErr:  false,
		},
		{
			name:     "non-existent email",
			email:    "nonexistent@example.com",
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetUserByEmail(ctx, tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantUser.ID, user.ID)
			assert.Equal(t, tt.wantUser.Email, user.Email)
		})
	}
}

func TestUserRepository_GetUsersChatList(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := mongodb.NewUserRepository(conn.DB)
	ctx := context.Background()

	testUser := &models.User{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "user",
	}
	err := repo.CreateUser(ctx, testUser)
	require.NoError(t, err)

	testChats := []bson.M{
		{
			"_id":   "chat-1",
			"title": "First Chat",
			"user":  testUser.ID,
		},
		{
			"_id":   "chat-2",
			"title": "Second Chat",
			"user":  testUser.ID,
		},
		{
			"_id":   "chat-3",
			"title": "Other User Chat",
			"user":  "other-user-id",
		},
	}

	for _, chat := range testChats {
		_, err := conn.DB.Collection("chats").InsertOne(ctx, chat)
		require.NoError(t, err)
	}

	var count int64
	count, err = conn.DB.Collection("chats").CountDocuments(ctx, bson.M{"user": testUser.ID})
	require.NoError(t, err)
	require.Equal(t, int64(2), count)

	tests := []struct {
		name      string
		userID    string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "user with chats",
			userID:    "test-user-id",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "user without chats",
			userID:    "non-existent-user",
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chats, err := repo.GetUsersChatList(ctx, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, chats, tt.wantCount)

			if tt.wantCount > 0 {
				for _, chat := range chats {
					assert.Equal(t, tt.userID, chat.User)
					assert.Contains(t, []string{"First Chat", "Second Chat"}, chat.ChatTitle)
					assert.NotEmpty(t, chat.ID)
				}
			}
		})
	}
}
