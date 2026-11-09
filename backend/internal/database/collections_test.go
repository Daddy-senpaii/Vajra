package database

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestEnsureUsersCollectionIsIdempotentAndEmpty(t *testing.T) {
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

	name := "vajra_init_test_" + time.Now().UTC().Format("20060102150405.000000000")
	db := client.Database(name)
	defer func() { _ = db.Drop(context.Background()) }()

	if err := EnsureUsersCollection(ctx, db); err != nil {
		t.Fatalf("EnsureUsersCollection(): %v", err)
	}
	if err := EnsureUsersCollection(ctx, db); err != nil {
		t.Fatalf("second EnsureUsersCollection(): %v", err)
	}

	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		t.Fatal(err)
	}
	if len(collections) != 1 || collections[0] != "users" {
		t.Fatalf("collections = %v, want only [users]", collections)
	}

	count, err := db.Collection("users").CountDocuments(ctx, bson.D{})
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("users document count = %d, want 0", count)
	}

	indexes, err := db.Collection("users").Indexes().List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer indexes.Close(ctx)
	if !indexes.Next(ctx) {
		t.Fatal("expected MongoDB's built-in _id index")
	}
	var index struct {
		Name string `bson:"name"`
	}
	if err := indexes.Decode(&index); err != nil {
		t.Fatal(err)
	}
	if index.Name != "_id_" || indexes.Next(ctx) {
		t.Fatalf("users collection should have only MongoDB's built-in _id index; first index is %q", index.Name)
	}
}
