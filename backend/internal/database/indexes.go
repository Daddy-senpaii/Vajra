package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type collectionIndex struct {
	collection string
	model      mongo.IndexModel
}

// EnsureIndexes is safe to run at each startup and keeps index definitions in source control.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	unique := options.Index().SetUnique(true)
	indexes := []collectionIndex{
		{"users", mongo.IndexModel{Keys: bson.D{{Key: "username_normalized", Value: 1}}, Options: unique}},
		{"users", mongo.IndexModel{Keys: bson.D{{Key: "email_normalized", Value: 1}}, Options: unique}},
		{"challenge_templates", mongo.IndexModel{Keys: bson.D{{Key: "slug", Value: 1}}, Options: unique}},
		{"enrollments", mongo.IndexModel{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "status", Value: 1}}}},
		{"progress_entries", mongo.IndexModel{Keys: bson.D{{Key: "enrollment_id", Value: 1}, {Key: "rule_id", Value: 1}, {Key: "local_date", Value: 1}}, Options: unique}},
		{"friend_requests", mongo.IndexModel{Keys: bson.D{{Key: "sender_id", Value: 1}, {Key: "recipient_id", Value: 1}, {Key: "status", Value: 1}}}},
		{"friendships", mongo.IndexModel{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "friend_id", Value: 1}}, Options: unique}},
		{"activity_events", mongo.IndexModel{Keys: bson.D{{Key: "recipient_id", Value: 1}, {Key: "created_at", Value: -1}}}},
		{"notifications", mongo.IndexModel{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}}}},
		{"friend_profile_views", mongo.IndexModel{Keys: bson.D{{Key: "viewer_id", Value: 1}, {Key: "profile_user_id", Value: 1}}, Options: unique}},
		{"running_entries", mongo.IndexModel{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "started_at", Value: -1}}}},
		{"reading_sessions", mongo.IndexModel{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "recorded_at", Value: -1}}}},
		{"workout_sessions", mongo.IndexModel{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "recorded_at", Value: -1}}}},
		{"digital_wellbeing_entries", mongo.IndexModel{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "local_date", Value: -1}}, Options: unique}},
	}
	for _, item := range indexes {
		if _, err := db.Collection(item.collection).Indexes().CreateOne(ctx, item.model); err != nil {
			return fmt.Errorf("index collection %q: %w", item.collection, err)
		}
	}
	return nil
}
