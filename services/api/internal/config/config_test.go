package config

import "testing"

func TestFromEnvDefaultsToPublicAdultFilterWithoutAPIKey(t *testing.T) {
	t.Setenv("CLEANBROWSING_API_KEY", "")
	t.Setenv("CLEANBROWSING_API_URL", "")
	t.Setenv("CLEANBROWSING_DOT_HOST", "")
	t.Setenv("CLEANBROWSING_DOH_URL", "")
	cfg, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv() error = %v", err)
	}
	if cfg.CleanBrowsingAPIKey != "" {
		t.Fatal("API key should be optional for public Adult Filter mode")
	}
	if cfg.CleanBrowsingDoTHost != "adult-filter-dns.cleanbrowsing.org" {
		t.Fatalf("DoT host = %q", cfg.CleanBrowsingDoTHost)
	}
}

func TestFromEnvRejectsInsecureRemoteProviderAPI(t *testing.T) {
	t.Setenv("CLEANBROWSING_API_URL", "http://provider.example/api")
	if _, err := FromEnv(); err == nil {
		t.Fatal("FromEnv() accepted an insecure remote API URL")
	}
}
