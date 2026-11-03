package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/domain"
)

type Config struct {
	ListenAddr           string
	CleanBrowsingAPIKey  string
	CleanBrowsingAPIURL  string
	CleanBrowsingProfile string
	CleanBrowsingDoTHost string
	CleanBrowsingDoHURL  string
	AdminToken           string
	AllowedOrigins       []string
	ReadHeaderTimeout    time.Duration
	ShutdownTimeout      time.Duration
}

func FromEnv() (Config, error) {
	apiURL := env("CLEANBROWSING_API_URL", "https://my.cleanbrowsing.org/api")
	adminToken := strings.TrimSpace(os.Getenv("VAJRA_ADMIN_TOKEN"))
	if adminToken != "" && len(adminToken) < 32 {
		return Config{}, fmt.Errorf("VAJRA_ADMIN_TOKEN must contain at least 32 characters")
	}
	parsed, err := url.Parse(apiURL)
	localHTTP := parsed != nil && parsed.Scheme == "http" && (parsed.Hostname() == "127.0.0.1" || parsed.Hostname() == "localhost")
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || (parsed.Scheme != "https" && !localHTTP) || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return Config{}, fmt.Errorf("CLEANBROWSING_API_URL must be an HTTPS URL")
	}
	dohURL := env("CLEANBROWSING_DOH_URL", "https://doh.cleanbrowsing.org/doh/adult-filter/")
	parsedDoH, err := url.Parse(dohURL)
	if err != nil || parsedDoH.Scheme != "https" || parsedDoH.Host == "" || parsedDoH.User != nil || parsedDoH.Fragment != "" {
		return Config{}, fmt.Errorf("CLEANBROWSING_DOH_URL must be an HTTPS URL")
	}
	dotHost := env("CLEANBROWSING_DOT_HOST", "adult-filter-dns.cleanbrowsing.org")
	if strings.ContainsAny(dotHost, "/:?#") {
		return Config{}, fmt.Errorf("CLEANBROWSING_DOT_HOST must be a hostname")
	}
	normalizedDoTHost, err := domain.Normalize(dotHost)
	if err != nil || normalizedDoTHost != strings.ToLower(strings.TrimSuffix(dotHost, ".")) {
		return Config{}, fmt.Errorf("CLEANBROWSING_DOT_HOST must be a valid normalized hostname")
	}
	dotHost = normalizedDoTHost

	return Config{
		ListenAddr:           env("VAJRA_API_LISTEN_ADDR", ":8080"),
		CleanBrowsingAPIKey:  strings.TrimSpace(os.Getenv("CLEANBROWSING_API_KEY")),
		CleanBrowsingAPIURL:  apiURL,
		CleanBrowsingProfile: env("CLEANBROWSING_PROFILE_NAME", "default"),
		CleanBrowsingDoTHost: dotHost,
		CleanBrowsingDoHURL:  dohURL,
		AdminToken:           adminToken,
		AllowedOrigins:       splitCSV(os.Getenv("VAJRA_API_ALLOWED_ORIGINS")),
		ReadHeaderTimeout:    5 * time.Second,
		ShutdownTimeout:      10 * time.Second,
	}, nil
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	var values []string
	for _, item := range strings.Split(value, ",") {
		if item = strings.TrimSpace(item); item != "" {
			values = append(values, item)
		}
	}
	return values
}
