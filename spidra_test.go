package spidra

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestClientInitialization(t *testing.T) {
	apiKey := "test_api_key_12345"
	client := New(apiKey)

	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Scrape == nil || client.Search == nil || client.Map == nil ||
		client.Parse == nil || client.Crawl == nil || client.Batch == nil ||
		client.Templates == nil || client.Presets == nil || client.Monitors == nil ||
		client.Account == nil || client.Logs == nil || client.Usage == nil {
		t.Fatal("all client resources must be initialized")
	}

	customClient := &http.Client{Timeout: 45 * time.Second}
	client2 := New(apiKey,
		WithBaseURL("https://custom.api.spidra.io"),
		WithTimeout(15*time.Second),
		WithHTTPClient(customClient),
	)

	if client2.Scrape.http.baseURL != "https://custom.api.spidra.io" {
		t.Fatalf("expected custom base URL, got: %s", client2.Scrape.http.baseURL)
	}
}

func TestPollOptionsDefaults(t *testing.T) {
	opts := PollOptions{}.withDefaults()
	if opts.Timeout != 120*time.Second {
		t.Errorf("expected default timeout 120s, got %v", opts.Timeout)
	}
	if opts.PollInterval != 3*time.Second {
		t.Errorf("expected default poll interval 3s, got %v", opts.PollInterval)
	}

	custom := PollOptions{
		Timeout:      60 * time.Second,
		PollInterval: 1 * time.Second,
	}.withDefaults()

	if custom.Timeout != 60*time.Second || custom.PollInterval != 1*time.Second {
		t.Errorf("expected custom poll options preserved")
	}
}

func TestMapErrorStatus(t *testing.T) {
	tests := []struct {
		status   int
		body     string
		checkErr func(error) bool
	}{
		{
			status: 401,
			body:   `{"error":"Invalid API key"}`,
			checkErr: func(err error) bool {
				var authErr *AuthenticationError
				return errors.As(err, &authErr) && authErr.StatusCode == 401
			},
		},
		{
			status: 402,
			body:   `{"error":"Insufficient credits for scrape"}`,
			checkErr: func(err error) bool {
				var credErr *InsufficientCreditsError
				return errors.As(err, &credErr) && credErr.StatusCode == 402
			},
		},
		{
			status: 403,
			body:   `{"error":"Access forbidden"}`,
			checkErr: func(err error) bool {
				var fErr *ForbiddenError
				return errors.As(err, &fErr) && fErr.StatusCode == 403
			},
		},
		{
			status: 404,
			body:   `{"error":"Job not found"}`,
			checkErr: func(err error) bool {
				var nErr *NotFoundError
				return errors.As(err, &nErr) && nErr.StatusCode == 404
			},
		},
		{
			status: 422,
			body:   `{"error":"Validation failed","details":["Invalid URL format","Limit must be > 0"]}`,
			checkErr: func(err error) bool {
				var vErr *ValidationError
				return errors.As(err, &vErr) && len(vErr.Errors) == 2
			},
		},
		{
			status: 429,
			body:   `{"error":"Rate limit exceeded"}`,
			checkErr: func(err error) bool {
				var rErr *RateLimitError
				return errors.As(err, &rErr) && rErr.StatusCode == 429
			},
		},
		{
			status: 500,
			body:   `{"error":"Internal server error"}`,
			checkErr: func(err error) bool {
				var sErr *ServerError
				return errors.As(err, &sErr) && sErr.StatusCode == 500
			},
		},
	}

	for _, tt := range tests {
		err := mapError(tt.status, []byte(tt.body))
		if !tt.checkErr(err) {
			t.Errorf("mapError(%d) failed check: %v", tt.status, err)
		}
	}
}

func TestWebhookVerification(t *testing.T) {
	payload := `{"event":"monitor.change","monitorId":"mon_123","hasChanges":true}`
	secret := "whsec_test_secret_key_abcdef123456"

	// Compute expected HMAC-SHA256 signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	computedHex := hex.EncodeToString(mac.Sum(nil))

	// Test with sha256= prefix
	sigWithPrefix := "sha256=" + computedHex
	if !VerifyWebhookString(payload, sigWithPrefix, secret) {
		t.Errorf("expected signature with prefix to be valid")
	}

	// Test without prefix
	if !VerifyWebhookString(payload, computedHex, secret) {
		t.Errorf("expected raw hex signature to be valid")
	}

	// Test byte slice variant
	if !VerifyWebhook([]byte(payload), sigWithPrefix, secret) {
		t.Errorf("expected byte slice verify to be valid")
	}

	// Test invalid signature
	if VerifyWebhookString(payload, "sha256=0000000000000000000000000000000000000000000000000000000000000000", secret) {
		t.Errorf("expected corrupted signature to be rejected")
	}

	// Test wrong secret
	if VerifyWebhookString(payload, sigWithPrefix, "wrong_secret") {
		t.Errorf("expected wrong secret to be rejected")
	}

	// Test empty parameters
	if VerifyWebhookString("", sigWithPrefix, secret) {
		t.Errorf("expected empty payload to be rejected")
	}
	if VerifyWebhookString(payload, "", secret) {
		t.Errorf("expected empty signature to be rejected")
	}
	if VerifyWebhookString(payload, sigWithPrefix, "") {
		t.Errorf("expected empty secret to be rejected")
	}
}

func TestJSONSerialization(t *testing.T) {
	// Test SearchParams
	searchParams := SearchParams{
		Query:   "latest AI news",
		Sources: []string{"web", "news"},
		Limit:   5,
		Scrape:  true,
		ScrapeOptions: &SearchScrapeOptions{
			Output: "markdown",
		},
	}
	b, err := json.Marshal(searchParams)
	if err != nil {
		t.Fatalf("failed to marshal SearchParams: %v", err)
	}
	var unmarshaled SearchParams
	if err := json.Unmarshal(b, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal SearchParams: %v", err)
	}
	if unmarshaled.Query != "latest AI news" || len(unmarshaled.Sources) != 2 || !unmarshaled.Scrape {
		t.Fatalf("unmarshaled SearchParams did not match: %+v", unmarshaled)
	}

	// Test MapParams
	mapParams := MapParams{
		URL:               "https://example.com",
		Search:            "pricing",
		IncludeSubdomains: true,
	}
	b, err = json.Marshal(mapParams)
	if err != nil {
		t.Fatalf("failed to marshal MapParams: %v", err)
	}
	var unmarshaledMap MapParams
	if err := json.Unmarshal(b, &unmarshaledMap); err != nil {
		t.Fatalf("failed to unmarshal MapParams: %v", err)
	}
	if unmarshaledMap.URL != "https://example.com" || unmarshaledMap.Search != "pricing" || !unmarshaledMap.IncludeSubdomains {
		t.Fatalf("unmarshaled MapParams did not match: %+v", unmarshaledMap)
	}

	// Test ScrapeParams with BrowserAction
	scrapeParams := ScrapeParams{
		URLs: []ScrapeURL{
			{URL: "https://example.com"},
		},
		Actions: []BrowserAction{
			{
				Type:     "click",
				Selector: "button.accept-cookies",
			},
			{
				Type:         "wait",
				Milliseconds: 1000,
			},
		},
	}
	b, err = json.Marshal(scrapeParams)
	if err != nil {
		t.Fatalf("failed to marshal ScrapeParams: %v", err)
	}
	var unmarshaledScrape ScrapeParams
	if err := json.Unmarshal(b, &unmarshaledScrape); err != nil {
		t.Fatalf("failed to unmarshal ScrapeParams: %v", err)
	}
	if len(unmarshaledScrape.Actions) != 2 || unmarshaledScrape.Actions[0].Type != "click" {
		t.Fatalf("unmarshaled ScrapeParams did not match: %+v", unmarshaledScrape)
	}
}

