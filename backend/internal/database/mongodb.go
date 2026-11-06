package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	Client *mongo.Client
	Name   string
}

func Connect(ctx context.Context, uri, name string) (*Database, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("create MongoDB client: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}
	database := client.Database(name)
	if err := EnsureIndexes(ctx, database); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("create MongoDB indexes: %w", err)
	}
	return &Database{Client: client, Name: name}, nil
}
