package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBlocklistEndpointIncludesETagAndConditionalResponse(t *testing.T) {
	handler := New(Options{Domains: []string{"second.example", "first.example"}}).Handler()
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/v1/blocklist", nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want 200", first.Code)
	}
	if !strings.Contains(first.Body.String(), `"count":2`) || !strings.Contains(first.Body.String(), `"first.example"`) {
		t.Fatalf("unexpected blocklist body: %s", first.Body.String())
	}

	request := httptest.NewRequest(http.MethodGet, "/v1/blocklist", nil)
	request.Header.Set("If-None-Match", first.Header().Get("ETag"))
	second := httptest.NewRecorder()
	handler.ServeHTTP(second, request)
	if second.Code != http.StatusNotModified {
		t.Fatalf("conditional status = %d, want 304", second.Code)
	}

	request = httptest.NewRequest(http.MethodGet, "/v1/blocklist", nil)
	request.Header.Set("If-None-Match", `"another-version", W/`+first.Header().Get("ETag"))
	second = httptest.NewRecorder()
	handler.ServeHTTP(second, request)
	if second.Code != http.StatusNotModified {
		t.Fatalf("weak conditional status = %d, want 304", second.Code)
	}
}

func TestAdminSyncRequiresBearerToken(t *testing.T) {
	handler := New(Options{AdminToken: "admin-secret"}).Handler()
	request := httptest.NewRequest(http.MethodPost, "/v1/admin/cleanbrowsing/sync", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated sync status = %d, want 401", response.Code)
	}

	request = httptest.NewRequest(http.MethodPost, "/v1/admin/cleanbrowsing/sync", nil)
	request.Header.Set("Authorization", "Bearer admin-secret")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("sync without provider status = %d, want 503", response.Code)
	}
}

func TestConfigDoesNotExposeProviderAPIKey(t *testing.T) {
	handler := New(Options{Config: Config{Provider: "CleanBrowsing", Filter: "adult", DoTHost: "adult.example"}}).Handler()
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/config", nil))
	if strings.Contains(response.Body.String(), "apikey") || strings.Contains(response.Body.String(), "secret") {
		t.Fatalf("config response contains a credential: %s", response.Body.String())
	}
}
