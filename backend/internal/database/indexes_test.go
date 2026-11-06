package database

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestEnsureIndexesIsRepeatable(t *testing.T) {
	uri := os.Getenv("VAYRA_TEST_MONGODB_URI")
	if uri == "" {
		t.Skip("set VAYRA_TEST_MONGODB_URI to run MongoDB integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	name := "vayra_index_test_" + time.Now().UTC().Format("20060102150405.000000000")
	db := client.Database(name)
	defer func() { _ = db.Drop(context.Background()) }()

	if err := EnsureIndexes(ctx, db); err != nil {
		t.Fatalf("first EnsureIndexes(): %v", err)
	}
	if err := EnsureIndexes(ctx, db); err != nil {
		t.Fatalf("second EnsureIndexes(): %v", err)
	}
	indexes, err := db.Collection("progress_entries").Indexes().List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer indexes.Close(ctx)
	if !indexes.Next(ctx) {
		t.Fatal("expected progress entry indexes")
	}
}
