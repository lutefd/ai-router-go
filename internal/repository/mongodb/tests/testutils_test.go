package mongodb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lutefd/ai-router-go/internal/database/mongodb"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*mongodb.MongoDBConnection, func()) {
	ctx := context.Background()

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForLog("Waiting for connections"),
		},
		Started: true,
	})
	require.NoError(t, err)

	port, err := mongoContainer.MappedPort(ctx, "27017")
	require.NoError(t, err)

	host, err := mongoContainer.Host(ctx)
	require.NoError(t, err)

	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	conn, err := mongodb.NewMongoDBConnection(ctx, mongoURI, "test_db")
	require.NoError(t, err)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		conn.Close(ctx)
		mongoContainer.Terminate(ctx)
	}

	return conn, cleanup
}
