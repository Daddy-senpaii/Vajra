package config

import "testing"

func TestLoadUsesEnvironmentAndDefaults(t *testing.T) {
	t.Setenv("VAYRA_LISTEN_ADDR", "127.0.0.1:9090")
	t.Setenv("VAYRA_MONGODB_URI", "mongodb://example.test:27017")
	t.Setenv("VAYRA_MONGODB_DATABASE", "testdb")

	cfg := Load()
	if cfg.ListenAddr != "127.0.0.1:9090" || cfg.MongoURI != "mongodb://example.test:27017" || cfg.MongoDatabase != "testdb" {
		t.Fatalf("Load() = %+v, want environment values", cfg)
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("VAYRA_LISTEN_ADDR", "")
	t.Setenv("VAYRA_MONGODB_URI", "")
	t.Setenv("VAYRA_MONGODB_DATABASE", "")

	cfg := Load()
	if cfg.ListenAddr != ":8080" || cfg.MongoURI != "mongodb://localhost:27017" || cfg.MongoDatabase != "vajra" {
		t.Fatalf("Load() = %+v, want defaults", cfg)
	}
}
