# Spidra Go SDK

The official Go SDK for [Spidra](https://spidra.io) — AI-native web scraping, multi-engine search, file & document parsing, deep site crawling, turnkey extraction templates, saved presets, and scheduled change monitors. All results come back as typed Go structs ready for production pipelines.

[![Go Reference](https://pkg.go.dev/badge/github.com/spidra-io/spidra-go.svg)](https://pkg.go.dev/github.com/spidra-io/spidra-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-00ADD8.svg?style=flat-square)](https://golang.org)

---

## Features

- **Web Scraping**: Fast raw content extraction (Markdown, HTML, text, screenshots) or structured AI extraction with guaranteed JSON schemas.
- **Browser Actions**: Automated headless browser interactions: click, type, wait for selectors, infinite scroll, JavaScript evaluation, element screenshots, and multi-step workflows.
- **Multi-Engine Search**: Web search across Google, Reddit, YouTube, TikTok, Hacker News, research papers, and developer issues with optional automatic markdown page scraping.
- **Domain Mapping**: Fast URL discovery via sitemaps and site graphs with subdomain and robots.txt compliance.
- **Document & Media Parsing**: Parse PDFs, Office documents (DOCX, XLSX), audio, and video files with AI or fast layout modes.
- **Deep Site Crawling**: Recursive crawling with natural language instructions, link depth limits, parallel workers, crawl diagnostic reports, and per-page retries.
- **Batch Scraping**: Concurrent scraping of hundreds or thousands of URLs with automatic chunking and cookie-based retry support.
- **Turnkey Templates**: 100+ ready-made extraction workflows for platforms like Google Maps, Amazon, flights, LinkedIn, and real estate.
- **Saved Presets**: Save and rerun common tool and template configurations with a single identifier.
- **Scheduled Monitors**: Watch pages for changes, evaluate CSS or AI conditions, and send webhook/email/Slack notifications.
- **Webhook Verification**: Secure constant-time HMAC-SHA256 signature verification helper.
- **Account & Logs**: Check credit balances, inspect proxy usage, and review detailed execution logs.

---

## Requirements

- Go >= 1.21
- **Zero external dependencies** — built entirely on the Go standard library (`net/http`, `encoding/json`, `crypto/hmac`, `context`, etc.).

---

## Installation

```bash
go get github.com/spidra-io/spidra-go
```

Get your API key at [app.spidra.io](https://app.spidra.io) under **Settings → API Keys**.

---

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	spidra "github.com/spidra-io/spidra-go"
)

func main() {
	client := spidra.New("spd_YOUR_API_KEY")

	// Scrape structured data with AI
	job, err := client.Scrape.Run(context.Background(), spidra.ScrapeParams{
		URLs:   []spidra.ScrapeURL{{URL: "https://news.ycombinator.com"}},
		Prompt: "Extract the top 5 articles with title, points, author, and URL",
		Output: "json",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Content:", job.Result.Content)
}
```

---

## Table of Contents

- [Client Initialization](#client-initialization)
- [Scraping](#scraping)
  - [Simple scrape](#simple-scrape)
  - [Structured AI extraction with JSON schema](#structured-ai-extraction-with-json-schema)
  - [Browser actions & interactions](#browser-actions--interactions)
  - [Geo-targeted scraping](#geo-targeted-scraping)
  - [Raw formats (Markdown, Screenshots)](#raw-formats-markdown-screenshots)
- [Search](#search)
  - [Basic web search](#basic-web-search)
  - [Search with automatic page scraping](#search-with-automatic-page-scraping)
  - [Specialized sources and domain filters](#specialized-sources-and-domain-filters)
- [Map](#map)
- [Document Parsing](#document-parsing)
- [Crawling](#crawling)
  - [Start and wait for a crawl](#start-and-wait-for-a-crawl)
  - [Crawl reports and page retries](#crawl-reports-and-page-retries)
- [Batch Scraping](#batch-scraping)
- [Ready-Made Templates](#ready-made-templates)
- [Saved Presets](#saved-presets)
- [Scheduled Monitors](#scheduled-monitors)
- [Account & Usage](#account--usage)
- [Execution Logs](#execution-logs)
- [Webhook Verification](#webhook-verification)
- [Error Handling](#error-handling)

---

## Client Initialization

```go
import (
	"net/http"
	"time"

	spidra "github.com/spidra-io/spidra-go"
)

// Default configuration
client := spidra.New("spd_YOUR_API_KEY")

// With custom options
client = spidra.New("spd_YOUR_API_KEY",
	spidra.WithBaseURL("https://api.spidra.io/api"),
	spidra.WithTimeout(60*time.Second),
	spidra.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
)
```

---

## Scraping

### Simple scrape

Scrape a URL and wait for the result:

```go
job, err := client.Scrape.RunURL(ctx, "https://example.com")
if err != nil {
	log.Fatal(err)
}
fmt.Println(job.Result.Content)
```

### Structured AI extraction with JSON schema

Enforce structured JSON output using a standard JSON Schema:

```go
schema := map[string]any{
	"type": "object",
	"required": []string{"company", "valuation"},
	"properties": map[string]any{
		"company":   map[string]any{"type": "string"},
		"valuation": map[string]any{"type": "string"},
	},
}

job, err := client.Scrape.Run(ctx, spidra.ScrapeParams{
	URLs:   []spidra.ScrapeURL{{URL: "https://example.com/company"}},
	Prompt: "Extract company details and latest valuation",
	Output: "json",
	Schema: schema,
})
if err != nil {
	log.Fatal(err)
}
fmt.Println(job.Result.Content)
```

### Browser actions & interactions

Execute user actions (click buttons, fill input fields, scroll, wait) before extraction:

```go
job, err := client.Scrape.Run(ctx, spidra.ScrapeParams{
	URLs: []spidra.ScrapeURL{
		{
			URL: "https://example.com/login-protected",
			Actions: []map[string]any{
				{"type": "click", "selector": "button#cookie-consent"},
				{"type": "wait", "milliseconds": 500},
				{"type": "scroll", "direction": "down", "amount": 1000},
			},
		},
	},
	Prompt: "Extract product catalog items",
})
```

### Geo-targeted scraping

Route requests through residential proxies in specific countries:

```go
job, err := client.Scrape.Run(ctx, spidra.ScrapeParams{
	URLs:         []spidra.ScrapeURL{{URL: "https://example.com/prices"}},
	UseProxy:     true,
	ProxyCountry: "US", // US, GB, DE, FR, JP, etc.
})
```

### Raw formats (Markdown, Screenshots)

Capture full-page screenshots or pure markdown:

```go
job, err := client.Scrape.Run(ctx, spidra.ScrapeParams{
	URLs:               []spidra.ScrapeURL{{URL: "https://example.com"}},
	Output:             "markdown",
	FullPageScreenshot: true,
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Markdown content:", job.Result.Content)
fmt.Println("Screenshot URL:", job.Result.Screenshots[0])
```

---

## Search

Query search engines with natural language, recency filters, and optional markdown scraping:

### Basic web search

```go
job, err := client.Search.RunQuery(ctx, "latest advances in quantum computing")
if err != nil {
	log.Fatal(err)
}

for _, item := range job.Result.Results {
	fmt.Printf("%s - %s\n", item.Title, item.URL)
}
```

### Search with automatic page scraping

Fetch search results and automatically scrape each target page into markdown:

```go
job, err := client.Search.Run(ctx, spidra.SearchParams{
	Query:   "best noise cancelling headphones 2026",
	Sources: []string{"web"},
	Limit:   5,
	Scrape:  true,
	ScrapeOptions: &spidra.SearchScrapeOptions{
		Output: "markdown",
	},
})
if err != nil {
	log.Fatal(err)
}

for _, item := range job.Result.Results {
	fmt.Println("URL:", item.URL)
	fmt.Println("Scraped Content Preview:", item.ScrapedContent)
}
```

### Specialized sources and domain filters

```go
job, err := client.Search.Run(ctx, spidra.SearchParams{
	Query:   "Go concurrency patterns",
	Sources: []string{"developer", "web"},
	Domains: []string{"github.com", "golang.org"},
	Recency: "past_year",
	Limit:   10,
})
```

---

## Map

Discover all reachable URLs on a website using site graphs and sitemaps:

```go
res, err := client.Map.RunURL(ctx, "https://docs.docker.com")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Discovered %d links\n", res.TotalLinks)
for _, link := range res.Links {
	fmt.Println(link.URL)
}

// Or filter links matching a keyword
filtered, err := client.Map.RunURLWithSearch(ctx, "https://docs.docker.com", "compose")
```

---

## Document Parsing

Parse PDFs, DOCX, XLSX, audio, and video files with layout extraction:

```go
job, err := client.Parse.RunURL(ctx, "https://example.com/annual_report.pdf")
if err != nil {
	log.Fatal(err)
}

fmt.Println("Extracted Markdown:\n", job.Result.Markdown)
for _, page := range job.Pages {
	fmt.Printf("Page %d: %d characters\n", page.PageNumber, len(page.Markdown))
}
```

### Estimate parsing costs

```go
estimate, err := client.Parse.EstimateURL(ctx, "https://example.com/large_document.pdf")
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Pages: %d, Estimated Credits: %.2f\n", estimate.PageCount, estimate.EstimatedCredits)
```

---

## Crawling

### Start and wait for a crawl

Recursively crawl an entire website with natural language guidance:

```go
job, err := client.Crawl.Run(ctx, spidra.CrawlParams{
	BaseURL:              "https://blog.rust-lang.org",
	CrawlInstruction:     "Only follow blog release notes and announcement links",
	TransformInstruction: "Extract the article title, author, date, and key takeaways",
	MaxPages:             25,
})
if err != nil {
	log.Fatal(err)
}

for _, page := range job.Result {
	fmt.Printf("%s: %s\n", page.Title, page.URL)
	fmt.Println("Data:", page.Data)
}
```

### Crawl reports and page retries

Inspect skipped URLs or retry specific failed pages:

```go
// Fetch crawl diagnostic report
report, err := client.Crawl.Report(ctx, job.JobID)
if err == nil {
	fmt.Printf("Skipped URLs: %d\n", len(report.Report.SkippedURLs))
}

// Retry a failed page
retryResp, err := client.Crawl.RetryPage(ctx, job.JobID, "page_123")

// Cancel an active crawl
cancelResp, err := client.Crawl.Cancel(ctx, job.JobID)
```

---

## Batch Scraping

Scrape multiple URLs concurrently with automatic rate limiting:

```go
urls := []string{
	"https://example.com/page1",
	"https://example.com/page2",
	"https://example.com/page3",
}

batch, err := client.Batch.RunURLs(ctx, urls)
if err != nil {
	log.Fatal(err)
}

for _, item := range batch.Items {
	fmt.Printf("[%s] %s -> %v\n", item.Status, item.URL, item.Result)
}

// Retry failed batch items with fresh session cookies
err = client.Batch.Retry(ctx, batch.BatchID, "session_id=abcdef")
```

---

## Ready-Made Templates

Execute turnkey scrapers for 100+ popular platforms:

```go
// List available templates
templates, err := client.Templates.List(ctx, map[string]string{"category": "e-commerce"})
for _, t := range templates.Templates {
	fmt.Println(t.Slug, "-", t.Name)
}

// Estimate template cost
est, err := client.Templates.Estimate(ctx, "amazon-product-reviews", map[string]any{
	"url": "https://amazon.com/dp/B08N5WRWNW",
})
fmt.Printf("Estimated credits: %.2f\n", est.EstimatedCredits)

// Run template and wait for extracted rows
run, err := client.Templates.Run(ctx, "amazon-product-reviews", map[string]any{
	"url":      "https://amazon.com/dp/B08N5WRWNW",
	"maxItems": 50,
})
if err != nil {
	log.Fatal(err)
}

for _, row := range run.Result {
	fmt.Println("Review:", row["title"], "-", row["rating"])
}
```

---

## Saved Presets

Save commonly used scraping or template configurations and re-execute them by ID:

```go
// Create a preset
preset, err := client.Presets.Create(ctx, spidra.PresetParams{
	Name: "Tech News Daily",
	Tool: "scrape",
	Config: map[string]any{
		"urls":   []map[string]any{{"url": "https://news.ycombinator.com"}},
		"prompt": "Extract top stories",
		"output": "json",
	},
})

// Run the preset with optional overrides
res, err := client.Presets.Run(ctx, preset.ID, nil, true)
fmt.Println("Result:", res.Result)

// List presets
list, err := client.Presets.List(ctx, nil)

// Delete preset
err = client.Presets.Delete(ctx, preset.ID)
```

---

## Scheduled Monitors

Monitor websites for changes on a recurring schedule with notifications:

```go
// Create a monitor
monitor, err := client.Monitors.Create(ctx, spidra.MonitorParams{
	Name:      "Competitor Pricing",
	URL:       "https://example.com/pricing",
	Schedule:  "0 * * * *", // hourly cron
	Selector:  ".price-table",
	Condition: "Notify if any plan price changes",
	Channels: map[string]any{
		"webhook": "https://myapi.com/webhooks/pricing",
		"email":   "alerts@mycompany.com",
	},
})

// Trigger an immediate check
checkResp, err := client.Monitors.CheckNow(ctx, monitor.ID)

// Pause and resume
paused, err := client.Monitors.Pause(ctx, monitor.ID)
resumed, err := client.Monitors.Resume(ctx, monitor.ID)

// Inspect historical checks
checks, err := client.Monitors.ListChecks(ctx, monitor.ID, nil)
for _, c := range checks.Checks {
	fmt.Printf("Checked at %s: hasChanges=%v\n", c.CheckedAt, c.HasChanges)
}
```

---

## Account & Usage

```go
// Account overview: plan, available credits, proxy bandwidth
account, err := client.Account.Get(ctx)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Plan: %s (Status: %s)\n", account.Plan.Name, account.Plan.Status)
fmt.Printf("Available Credits: %.2f / %.2f\n", account.Credits.Available, account.Credits.Total)
fmt.Printf("Residential Proxy: %.2f MB available\n", account.ProxyBandwidth.AvailableMB)

// Credit usage stats over time ("7d", "30d", "weekly")
usageRows, err := client.Account.Usage(ctx, "30d")
for _, row := range usageRows {
	fmt.Printf("%s: %d requests, %.2f credits\n", row.Date, row.Requests, row.Credits)
}

// Tool credit breakdown ("scrape", "crawl", "search", "parse", "map")
breakdown, err := client.Account.Breakdown(ctx, "7d", "", "")
for _, b := range breakdown {
	fmt.Printf("[%s] %s: %.2f credits (%d requests)\n", b.Date, b.Tool, b.Credits, b.Requests)
}
```

---

## Execution Logs

Retrieve audit and debug logs for all operations:

```go
// Scrape logs
scrapeLogs, err := client.Logs.Scrape(ctx, map[string]string{"limit": "10"})

// Search logs
searchLogs, err := client.Logs.Search(ctx, map[string]string{"limit": "10"})

// Map logs
mapLogs, err := client.Logs.Map(ctx, map[string]string{"limit": "10"})

// Parse logs
parseLogs, err := client.Logs.Parse(ctx, map[string]string{"limit": "10"})
```

---

## Webhook Verification

Verify incoming webhook requests from Spidra using constant-time HMAC-SHA256:

```go
package main

import (
	"io"
	"net/http"

	spidra "github.com/spidra-io/spidra-go"
)

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sigHeader := r.Header.Get("X-Spidra-Signature")
	secret := "whsec_YOUR_WEBHOOK_SECRET"

	if !spidra.VerifyWebhook(body, sigHeader, secret) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Signature is valid — process event
	w.WriteHeader(http.StatusOK)
}
```

---

## Error Handling

All SDK errors map to specific typed Go error structs that implement `error`:

```go
import "errors"

job, err := client.Scrape.Run(ctx, params)
if err != nil {
	var authErr *spidra.AuthenticationError
	var credErr *spidra.InsufficientCreditsError
	var valErr *spidra.ValidationError
	var rateErr *spidra.RateLimitError
	var jobErr *spidra.JobFailedError
	var timeoutErr *spidra.TimeoutError
	var spidraErr *spidra.SpidraError

	switch {
	case errors.As(err, &authErr):
		log.Printf("Authentication failed: %s (Status: %d)", authErr.Message, authErr.StatusCode)
	case errors.As(err, &credErr):
		log.Printf("Out of credits: %s", credErr.Message)
	case errors.As(err, &valErr):
		log.Printf("Validation error: %s (Details: %v)", valErr.Message, valErr.Errors)
	case errors.As(err, &rateErr):
		log.Printf("Rate limited. Retry after %d seconds", rateErr.RetryAfter)
	case errors.As(err, &jobErr):
		log.Printf("Job %s failed with status %s: %s", jobErr.JobID, jobErr.Status, jobErr.Message)
	case errors.As(err, &timeoutErr):
		log.Printf("Job %s timed out after %d seconds", timeoutErr.JobID, timeoutErr.TimeoutSeconds)
	case errors.As(err, &spidraErr):
		log.Printf("Spidra API error %d: %s", spidraErr.StatusCode, spidraErr.Message)
	default:
		log.Printf("Unexpected error: %v", err)
	}
}
```

---

## License

This SDK is distributed under the [MIT License](LICENSE).
