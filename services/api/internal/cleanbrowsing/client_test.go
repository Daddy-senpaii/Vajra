package cleanbrowsing

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestSyncDomainsAddsOnlyMissingEntries(t *testing.T) {
	var added []string
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.Query().Get("apikey"); got != "server-secret" {
			t.Errorf("API key query = %q, want configured secret", got)
		}
		body := ""
		switch r.URL.Query().Get("action") {
		case "profiles/list":
			body = `{"status":"success","profiles":["default"]}`
		case "blocklist/list":
			if got := r.URL.Query().Get("profile_name"); got != "default" {
				t.Errorf("profile_name = %q, want default", got)
			}
			body = `{"status":"success","blocklist":{"default":["ALREADY.example."]}}`
		case "blocklist/add":
			added = append(added, r.URL.Query().Get("domain_name"))
			body = `{"status":"success"}`
		default:
			t.Errorf("unexpected action %q", r.URL.Query().Get("action"))
			return nil, io.ErrUnexpectedEOF
		}
		return jsonResponse(http.StatusOK, body), nil
	})

	client := NewClient("https://cleanbrowsing.example/api", "server-secret", "default", &http.Client{Transport: transport})
	result, err := client.SyncDomains(context.Background(), []string{"already.example", "new.example", "new.example", "bad..domain"})
	if err != nil {
		t.Fatalf("SyncDomains() error = %v", err)
	}
	if result.Added != 1 || result.Skipped != 2 || !reflect.DeepEqual(result.Failed, []string{"bad..domain"}) {
		t.Fatalf("SyncDomains() result = %#v", result)
	}
	if want := []string{"new.example"}; !reflect.DeepEqual(added, want) {
		t.Fatalf("added domains = %v, want %v", added, want)
	}
}

func TestRequestErrorDoesNotExposeAPIKey(t *testing.T) {
	transport := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusInternalServerError, "failure"), nil
	})
	client := NewClient("https://cleanbrowsing.example/api", "secret-do-not-log", "default", &http.Client{Transport: transport})
	err := client.call(context.Background(), "profiles/list", url.Values{})
	if err == nil {
		t.Fatal("call() error = nil, want HTTP error")
	}
	if got := err.Error(); got == "" || containsSecret(got, "secret-do-not-log") {
		t.Fatalf("error message leaked secret or was empty: %q", got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func containsSecret(value, secret string) bool {
	return len(secret) > 0 && len(value) >= len(secret) && valueContains(value, secret)
}

func valueContains(value, part string) bool {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return true
		}
	}
	return false
}
