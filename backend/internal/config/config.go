package config

import "os"

type Config struct {
	ListenAddr    string
	MongoURI      string
	MongoDatabase string
}

func Load() Config {
	return Config{
		ListenAddr:    envOr("VAYRA_LISTEN_ADDR", ":8080"),
		MongoURI:      envOr("VAYRA_MONGODB_URI", "mongodb://localhost:27017"),
		MongoDatabase: envOr("VAYRA_MONGODB_DATABASE", "vajra"),
	}
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
