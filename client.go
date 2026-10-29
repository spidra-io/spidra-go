package spidra

import (
	"context"
	"net/http"
	"time"
)

// Client is the entry point for all Spidra API resources.
type Client struct {
	Scrape    *ScrapeResource
	Search    *SearchResource
	Map       *MapResource
	Parse     *ParseResource
	Crawl     *CrawlResource
	Batch     *BatchResource
	Templates *TemplatesResource
	Presets   *PresetsResource
	Monitors  *MonitorsResource
	Account   *AccountResource
	Logs      *LogsResource
	Usage     *UsageResource
}

// Option configures the Client.
type Option func(*httpClient)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(h *httpClient) { h.baseURL = url }
}

// WithTimeout sets the HTTP client timeout (default: 30s).
func WithTimeout(d time.Duration) Option {
	return func(h *httpClient) { h.http.Timeout = d }
}

// WithHTTPClient allows supplying a custom *http.Client.
func WithHTTPClient(client *http.Client) Option {
	return func(h *httpClient) {
		if client != nil {
			h.http = client
		}
	}
}

// New creates a new Spidra client authenticated with the given API key.
func New(apiKey string, opts ...Option) *Client {
	h := &httpClient{
		apiKey:  apiKey,
		baseURL: "https://api.spidra.io/api",
		http:    &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(h)
	}
	return &Client{
		Scrape:    &ScrapeResource{h},
		Search:    &SearchResource{h},
		Map:       &MapResource{h},
		Parse:     &ParseResource{h},
		Crawl:     &CrawlResource{h},
		Batch:     &BatchResource{h},
		Templates: &TemplatesResource{h},
		Presets:   &PresetsResource{h},
		Monitors:  &MonitorsResource{h},
		Account:   &AccountResource{h},
		Logs:      &LogsResource{h},
		Usage:     &UsageResource{h},
	}
}

// PollOptions controls how Run() methods poll for job completion.
type PollOptions struct {
	// Timeout is the maximum time to wait. Default: 120s.
	Timeout time.Duration
	// PollInterval is the time between status checks. Default: 3s.
	PollInterval time.Duration
}

func (o PollOptions) withDefaults() PollOptions {
	if o.Timeout == 0 {
		o.Timeout = 120 * time.Second
	}
	if o.PollInterval == 0 {
		o.PollInterval = 3 * time.Second
	}
	return o
}

// QuickRunTemplate is a top-level convenience method to execute a turnkey template and wait for completion.
func (c *Client) QuickRunTemplate(ctx context.Context, slug string, inputs map[string]any, opts ...PollOptions) (*TemplateRun, error) {
	return c.Templates.Run(ctx, slug, inputs, opts...)
}

// QuickRunPreset is a top-level convenience method to execute a saved preset.
func (c *Client) QuickRunPreset(ctx context.Context, presetID string, overrides map[string]any, wait bool) (*RunPresetResponse, error) {
	return c.Presets.Run(ctx, presetID, overrides, wait)
}

// GetAccount is a top-level convenience method to retrieve account details.
func (c *Client) GetAccount(ctx context.Context) (*AccountSummary, error) {
	return c.Account.Get(ctx)
}
