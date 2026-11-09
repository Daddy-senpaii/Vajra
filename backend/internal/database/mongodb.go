package database

import (
	"context"
	"errors"
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
	if err := EnsureUsersCollection(ctx, database); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("create users collection: %w", err)
	}
	return &Database{Client: client, Name: name}, nil
}

// EnsureUsersCollection creates the initially empty authentication collection if needed.
func EnsureUsersCollection(ctx context.Context, db *mongo.Database) error {
	if err := db.CreateCollection(ctx, "users"); err != nil {
		var commandErr mongo.CommandError
		if errors.As(err, &commandErr) && commandErr.Code == 48 { // NamespaceExists
			return nil
		}
		return err
	}
	return nil
}
