package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConnection struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDBConnection(ctx context.Context, uri string, dbName string) (*MongoDBConnection, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(dbName)
	return &MongoDBConnection{
		Client: client,
		DB:     db,
	}, nil
}

func (m *MongoDBConnection) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
