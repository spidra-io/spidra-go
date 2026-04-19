package spidra

import (
	"net/http"
	"time"
)

// Client is the entry point for all Spidra API resources.
type Client struct {
	Scrape *ScrapeResource
	Batch  *BatchResource
	Crawl  *CrawlResource
	Logs   *LogsResource
	Usage  *UsageResource
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
		Scrape: &ScrapeResource{h},
		Batch:  &BatchResource{h},
		Crawl:  &CrawlResource{h},
		Logs:   &LogsResource{h},
		Usage:  &UsageResource{h},
	}
}

// PollOptions controls how Run() polls for job completion.
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
