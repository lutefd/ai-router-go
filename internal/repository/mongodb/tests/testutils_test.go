package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/lutefd/ai-router-go/internal/database/mongodb"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *mongodb.MongoDBConnection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := mongodb.NewMongoDBConnection(ctx, "mongodb://localhost:27017", "test_db")
	require.NoError(t, err)

	err = conn.DB.Collection("users").Drop(ctx)
	require.NoError(t, err)
	err = conn.DB.Collection("chats").Drop(ctx)
	require.NoError(t, err)

	return conn
}
