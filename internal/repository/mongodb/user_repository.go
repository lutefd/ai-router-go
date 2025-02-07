package mongodb

import (
	"context"
	"fmt"

	"github.com/lutefd/ai-router-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	result, err := r.db.Collection("users").ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID string) error {
	result, err := r.db.Collection("users").DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := r.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("error decoding users: %w", err)
	}
	return users, nil
}

func (r *UserRepository) GetUsersChatList(ctx context.Context, userID string) ([]*models.UserChat, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"user": userID},
		},
		{
			"$project": bson.M{
				"_id":        0,
				"id":         "$_id",
				"user":       "$user",
				"chat_title": "$title",
			},
		},
	}

	cursor, err := r.db.Collection("chats").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error fetching user chats: %w", err)
	}
	defer cursor.Close(ctx)

	var chats []*models.UserChat
	if err = cursor.All(ctx, &chats); err != nil {
		return nil, fmt.Errorf("error decoding user chats: %w", err)
	}

	var rawDocs []bson.M
	cursor, _ = r.db.Collection("chats").Aggregate(ctx, pipeline)
	cursor.All(ctx, &rawDocs)
	fmt.Printf("Raw documents: %+v\n", rawDocs)

	return chats, nil
}
