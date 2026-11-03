package cleanbrowsing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/blocklist"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/domain"
)

const responseLimit = 2 << 20

type Client struct {
	baseURL    string
	apiKey     string
	profile    string
	httpClient *http.Client
}

type SyncResult struct {
	Profile string   `json:"profile"`
	Added   int      `json:"added"`
	Skipped int      `json:"skipped"`
	Failed  []string `json:"failed,omitempty"`
}

func NewClient(baseURL, apiKey, profile string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	clientCopy := *httpClient
	clientCopy.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &Client{baseURL: strings.TrimRight(baseURL, "?&"), apiKey: apiKey, profile: profile, httpClient: &clientCopy}
}

func (c *Client) SyncDomains(ctx context.Context, domains []string) (SyncResult, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return SyncResult{}, errors.New("CleanBrowsing API key is not configured")
	}
	result := SyncResult{Profile: c.profile}

	profiles, err := c.listProfiles(ctx)
	if err != nil {
		return result, fmt.Errorf("list CleanBrowsing profiles: %w", err)
	}
	if !containsFold(profiles, c.profile) {
		return result, fmt.Errorf("CleanBrowsing profile %q does not exist; create/configure it in the CleanBrowsing dashboard first", c.profile)
	}

	existing, err := c.listBlockedDomains(ctx)
	if err != nil {
		return result, fmt.Errorf("list CleanBrowsing blocklist: %w", err)
	}
	for _, rawDomain := range domains {
		normalized, normalizeErr := domain.Normalize(rawDomain)
		if normalizeErr != nil {
			result.Failed = append(result.Failed, rawDomain)
			continue
		}
		if existing[normalized] {
			result.Skipped++
			continue
		}
		params := url.Values{"domain_name": {normalized}, "profile_name": {c.profile}}
		if err := c.call(ctx, "blocklist/add", params); err != nil {
			result.Failed = append(result.Failed, normalized)
			continue
		}
		result.Added++
		existing[normalized] = true
	}
	sort.Strings(result.Failed)
	return result, nil
}

func (c *Client) listProfiles(ctx context.Context) ([]string, error) {
	data, err := c.request(ctx, "profiles/list", nil)
	if err != nil {
		return nil, err
	}
	var response struct {
		Profiles []string `json:"profiles"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("decode profiles response: %w", err)
	}
	return response.Profiles, nil
}

func (c *Client) listBlockedDomains(ctx context.Context) (map[string]bool, error) {
	data, err := c.request(ctx, "blocklist/list", url.Values{"profile_name": {c.profile}})
	if err != nil {
		return nil, err
	}
	var response struct {
		Blocklist map[string][]string `json:"blocklist"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("decode blocklist response: %w", err)
	}
	blocked := make(map[string]bool)
	for _, domains := range response.Blocklist {
		for _, rawDomain := range domains {
			if normalized, err := domain.Normalize(rawDomain); err == nil {
				blocked[normalized] = true
			}
		}
	}
	return blocked, nil
}

func (c *Client) call(ctx context.Context, action string, params url.Values) error {
	_, err := c.request(ctx, action, params)
	return err
}

func (c *Client) request(ctx context.Context, action string, params url.Values) ([]byte, error) {
	query := url.Values{}
	for key, values := range params {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	query.Set("apikey", c.apiKey)
	query.Set("action", action)
	endpoint := c.baseURL
	if strings.Contains(endpoint, "?") {
		endpoint += "&" + query.Encode()
	} else {
		endpoint += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create CleanBrowsing request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Vajra-API/0.1")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Do not wrap the URL-bearing transport error: the API key is in its query.
		return nil, errors.New("CleanBrowsing request failed")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, responseLimit+1))
	if err != nil {
		return nil, errors.New("read CleanBrowsing response failed")
	}
	if len(body) > responseLimit {
		return nil, errors.New("CleanBrowsing response exceeded size limit")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("CleanBrowsing returned HTTP %d", resp.StatusCode)
	}
	var status struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, errors.New("CleanBrowsing returned invalid JSON")
	}
	if !strings.EqualFold(status.Status, "success") {
		if status.Reason == "" {
			return nil, errors.New("CleanBrowsing API request failed")
		}
		if c.apiKey != "" {
			status.Reason = strings.ReplaceAll(status.Reason, c.apiKey, "[redacted]")
		}
		return nil, fmt.Errorf("CleanBrowsing API request failed: %s", status.Reason)
	}
	return body, nil
}

func containsFold(values []string, wanted string) bool {
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), wanted) {
			return true
		}
	}
	return false
}

// CuratedDomains returns the project-supplied domains used for provider sync.
func CuratedDomains() ([]string, error) {
	list, _, err := blocklist.NewCurated()
	if err != nil {
		return nil, err
	}
	return list.Domains(), nil
}
