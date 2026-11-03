package httpapi

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Daddy-senpaii/Vajra/services/api/internal/cleanbrowsing"
)

type Config struct {
	DoTHost  string `json:"dot_hostname"`
	DoHURL   string `json:"doh_url"`
	Filter   string `json:"filter"`
	Provider string `json:"provider"`
	Profile  string `json:"profile,omitempty"`
}

type Server struct {
	domains       []string
	version       string
	config        Config
	adminToken    string
	allowedOrigin map[string]bool
	cleanBrowsing *cleanbrowsing.Client
}

type Options struct {
	Domains        []string
	Config         Config
	AdminToken     string
	AllowedOrigins []string
	CleanBrowsing  *cleanbrowsing.Client
}

func New(options Options) *Server {
	domains := append([]string(nil), options.Domains...)
	sort.Strings(domains)
	data := strings.Join(domains, "\n")
	sum := sha256.Sum256([]byte(data))
	origins := make(map[string]bool, len(options.AllowedOrigins))
	for _, origin := range options.AllowedOrigins {
		origins[origin] = true
	}
	return &Server{
		domains:       domains,
		version:       hex.EncodeToString(sum[:]),
		config:        options.Config,
		adminToken:    options.AdminToken,
		allowedOrigin: origins,
		cleanBrowsing: options.CleanBrowsing,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /v1/config", s.getConfig)
	mux.HandleFunc("GET /v1/blocklist", s.getBlocklist)
	mux.HandleFunc("POST /v1/admin/cleanbrowsing/sync", s.syncCleanBrowsing)
	return s.cors(mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=300")
	writeJSON(w, http.StatusOK, map[string]any{
		"api_version":       "v1",
		"blocklist_version": s.version,
		"dns_filter":        s.config,
	})
}

func (s *Server) getBlocklist(w http.ResponseWriter, r *http.Request) {
	etag := `"` + s.version + `"`
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=86400")
	w.Header().Set("X-Blocklist-Version", s.version)
	if matchesETag(r.Header.Values("If-None-Match"), etag) {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"version": s.version,
		"source":  "vajra-curated",
		"count":   len(s.domains),
		"domains": s.domains,
	})
}

func matchesETag(values []string, etag string) bool {
	for _, value := range values {
		for _, candidate := range strings.Split(value, ",") {
			candidate = strings.TrimSpace(candidate)
			if candidate == "*" || candidate == etag || strings.TrimPrefix(candidate, "W/") == etag {
				return true
			}
		}
	}
	return false
}

func (s *Server) syncCleanBrowsing(w http.ResponseWriter, r *http.Request) {
	if s.adminToken == "" || !validBearer(r.Header.Get("Authorization"), s.adminToken) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if s.cleanBrowsing == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "CleanBrowsing API is not configured"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()
	result, err := s.cleanBrowsing.SyncDomains(ctx, s.domains)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	status := http.StatusOK
	if len(result.Failed) > 0 {
		status = http.StatusMultiStatus
	}
	writeJSON(w, status, result)
}

func (s *Server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && s.allowedOrigin[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, If-None-Match, Content-Type")
			w.Header().Set("Access-Control-Expose-Headers", "ETag, X-Blocklist-Version")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func validBearer(header, token string) bool {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return false
	}
	provided := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if len(provided) != len(token) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(provided), []byte(token)) == 1
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
