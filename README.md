# Spidra Go SDK

The official Go SDK for [Spidra](https://spidra.io) — scrape pages, run browser actions, batch-process URLs, and crawl entire sites. All results come back as typed structs ready to feed into your pipelines or store directly.

## Requirements

- Go >= 1.21
- No external dependencies — uses the standard library only

## Installation

```bash
go get github.com/spidra-io/spidra-go
```

Get your API key at [app.spidra.io](https://app.spidra.io) under **Settings → API Keys**.

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

    job, err := client.Scrape.Run(context.Background(), spidra.ScrapeParams{
        URLs:   []spidra.ScrapeURL{{URL: "https://news.ycombinator.com"}},
        Prompt: "List the top 5 stories with title, points, and comment count",
        Output: "json",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(job.Result.Content)
}
```

## Table of Contents

- [Scraping](#scraping)
  - [Basic scrape](#basic-scrape)
  - [Structured output with JSON schema](#structured-output-with-json-schema)
  - [Geo-targeted scraping](#geo-targeted-scraping)
  - [Authenticated pages](#authenticated-pages)
  - [Browser actions](#browser-actions)
  - [forEach: process every element on a page](#foreach-process-every-element-on-a-page)
  - [Manual job control](#manual-job-control)
- [Batch Scraping](#batch-scraping)
- [Crawling](#crawling)
- [Logs](#logs)
- [Usage Statistics](#usage-statistics)
- [Error Handling](#error-handling)

## Scraping

All scrape jobs run asynchronously. `Run()` submits a job and polls until it finishes. If you need more control, use `Submit()` and `Get()` directly.

Up to 3 URLs can be passed per request and they are processed in parallel.

### Basic scrape

```go
job, err := client.Scrape.Run(ctx, spidra.ScrapeParams{
    URLs:   []spidra.ScrapeURL{{URL: "https://example.com/pricing"}},
    Prompt: "Extract all pricing plans with name, price, and included features",
    Output: "json",
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(job.Result.Content)
// map[plans:[map[name:Starter price:$9/mo features:[...]] ...]]
```

Output can be `"json"` (default) or `"markdown"`. For markdown, `Content` is a string. For JSON, `Content` is a parsed object or array.

### Structured output with JSON schema

When you need a guaranteed shape, pass a `Schema`. The API will enforce the structure and return `null` for any missing fields rather than hallucinating values.

```go
job, err := client.Scrape.Run(ctx, spidra.ScrapeParams{
    URLs:   []spidra.ScrapeURL{{URL: "https://jobs.example.com/senior-engineer"}},
    Prompt: "Extract the job listing details",
    Output: "json",
    Schema: map[string]any{
        "type":     "object",
        "required": []string{"title", "company", "remote"},
        "properties": map[string]any{
            "title":      map[string]any{"type": "string"},
            "company":    map[string]any{"type": "string"},
            "remote":     map[string]any{"type": []any{"boolean", "null"}},
            "salary_min": map[string]any{"type": []any{"number", "null"}},
            "salary_max": map[string]any{"type": []any{"number", "null"}},
            "skills":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
        },
    },
})
```

### Geo-targeted scraping

Pass `UseProxy: true` and a `ProxyCountry` code to route the request through a specific country. Useful for geo-restricted content or localized pricing.

```go
job, err := client.Scrape.Run(ctx, spidra.ScrapeParams{
    URLs:         []spidra.ScrapeURL{{URL: "https://www.amazon.de/gp/bestsellers"}},
    Prompt:       "List the top 10 products with name and price",
    UseProxy:     true,
    ProxyCountry: "de",
})
```

Supported country codes include: `us`, `gb`, `de`, `fr`, `jp`, `au`, `ca`, `br`, `in`, `nl`, `sg`, `es`, `it`, `mx`, and [40+ more](https://docs.spidra.io/features/stealth-mode#country-targeting). Use `"global"` or `"eu"` for regional routing.

### Authenticated pages

Pass cookies as a string to scrape pages that require a login session.

```go
job, err := client.Scrape.Run(ctx, spidra.ScrapeParams{
    URLs:    []spidra.ScrapeURL{{URL: "https://app.example.com/dashboard"}},
    Prompt:  "Extract the monthly revenue and active user count",
    Cookies: "session=abc123; auth_token=xyz789",
})
```

### Browser actions

Actions let you interact with the page before the scrape runs. They execute in order, and the scrape happens after all actions complete.

```go
job, err := client.Scrape.Run(ctx, spidra.ScrapeParams{
    URLs: []spidra.ScrapeURL{{
        URL: "https://example.com/products",
        Actions: []map[string]any{
            {"type": "click", "selector": "#accept-cookies"},
            {"type": "wait", "duration": 1000},
            {"type": "scroll", "to": "80%"},
        },
    }},
    Prompt: "Extract all product names and prices",
})
```

**Available actions:**

| Action    | Required fields       | Description |
|-----------|-----------------------|-------------|
| `click`   | `selector` or `value` | Click a button, link, or any element |
| `type`    | `selector`, `value`   | Type text into an input or textarea |
| `check`   | `selector` or `value` | Check a checkbox |
| `uncheck` | `selector` or `value` | Uncheck a checkbox |
| `wait`    | `duration` (ms)       | Pause for a set number of milliseconds |
| `scroll`  | `to` (0–100%)         | Scroll the page to a percentage of its height |
| `forEach` | `observe`             | Loop over every matched element and process each one |

For `selector`, use a CSS selector. For `value`, use a plain English description and Spidra will locate the element using AI.

```go
// CSS selector
{"type": "click", "selector": "button[data-testid='submit']"}

// Plain English — AI finds the element
{"type": "click", "value": "Accept all cookies button"}

// Type into a field
{"type": "type", "selector": "input[name='q']", "value": "wireless headphones"}

// Wait for content to load
{"type": "wait", "duration": 2000}

// Scroll to bottom
{"type": "scroll", "to": "100%"}
```

### forEach: process every element on a page

`forEach` finds a set of elements on the page and processes each one individually. It is the right tool when you need to collect data from a list of items, paginate through multiple pages, or click into each item's detail page.

> You don't need `forEach` if the data fits on a single page and is short — a plain `Prompt` is simpler and works just as well.

**Use forEach when:**
- The list spans multiple pages and you need `pagination`
- You need to click into each item's detail page (`navigate` mode)
- You have 20+ items and want per-item AI extraction to stay consistent (`itemPrompt`)

#### inline mode

Read each element's content directly without navigating. Best for product cards, search results, table rows.

```go
job, err := client.Scrape.Run(ctx, spidra.ScrapeParams{
    URLs: []spidra.ScrapeURL{{
        URL: "https://books.toscrape.com/catalogue/category/books/mystery_3/index.html",
        Actions: []map[string]any{{
            "type":            "forEach",
            "observe":         "Find all book cards in the product grid",
            "mode":            "inline",
            "captureSelector": "article.product_pod",
            "maxItems":        20,
            "itemPrompt":      "Extract title, price, and star rating. Return as JSON: {title, price, star_rating}",
        }},
    }},
    Prompt: "Return a clean JSON array of all books",
    Output: "json",
})
```

#### navigate mode

Follow each element's link to its destination page and capture content there. Best for product listings where the full detail is only on the individual page.

```go
map[string]any{
    "type":            "forEach",
    "observe":         "Find all book title links in the product grid",
    "mode":            "navigate",
    "captureSelector": "article.product_page",
    "maxItems":        10,
    "waitAfterClick":  800,
    "itemPrompt":      "Extract title, price, star rating, and availability. Return as JSON.",
}
```

#### click mode

Click each element, capture the content that appears (a modal, drawer, or expanded section), then move on. Best for hotel room cards, FAQ accordions, or any UI where clicking reveals hidden content.

```go
map[string]any{
    "type":            "forEach",
    "observe":         "Find all room type cards",
    "mode":            "click",
    "captureSelector": "[role='dialog']",
    "maxItems":        8,
    "waitAfterClick":  1200,
    "itemPrompt":      "Extract room name, bed type, price per night, and amenities. Return as JSON.",
}
```

#### Pagination

After processing all elements on the current page, follow the next-page link and continue collecting.

```go
map[string]any{
    "type":     "forEach",
    "observe":  "Find all book title links",
    "mode":     "navigate",
    "maxItems": 40,
    "pagination": map[string]any{
        "nextSelector": "li.next > a",
        "maxPages":     3, // 3 additional pages beyond the first
    },
}
```

`maxItems` applies across all pages combined. The loop stops when you hit `maxItems`, run out of elements, or reach `maxPages`.

#### Per-element actions

Run additional browser actions on each item after navigating or clicking into it, before the content is captured.

```go
map[string]any{
    "type":            "forEach",
    "observe":         "Find all book title links",
    "mode":            "navigate",
    "captureSelector": "article.product_page",
    "maxItems":        5,
    "waitAfterClick":  1000,
    "actions": []map[string]any{
        {"type": "scroll", "to": "50%"},
    },
    "itemPrompt": "Extract title, price, and full description. Return as JSON.",
}
```

### Manual job control

Use `Submit()` and `Get()` when you want to manage polling yourself, or fire-and-forget and check back later.

```go
// Submit and get the jobId immediately
queued, err := client.Scrape.Submit(ctx, spidra.ScrapeParams{
    URLs:   []spidra.ScrapeURL{{URL: "https://example.com/listings"}},
    Prompt: "Extract all property listings",
    Output: "json",
})
if err != nil {
    log.Fatal(err)
}

jobID := queued.JobID

// Check status at any point
result, err := client.Scrape.Get(ctx, jobID)
if err != nil {
    log.Fatal(err)
}

if result.Status == "completed" {
    fmt.Println(result.Result.Content)
} else if result.Status == "failed" {
    fmt.Println("failed:", result.Error)
}
```

Job statuses: `waiting`, `active`, `completed`, `failed`.

**Custom timeout and poll interval:**

```go
job, err := client.Scrape.Run(ctx, params, spidra.PollOptions{
    Timeout:      3 * time.Minute, // wait up to 3 minutes
    PollInterval: 5 * time.Second, // check every 5 seconds
})
```

## Batch Scraping

Submit up to 50 URLs in a single request. All URLs are processed in parallel. Each URL is a plain string.

```go
batch, err := client.Batch.Run(ctx, spidra.BatchParams{
    URLs: []string{
        "https://shop.example.com/product/1",
        "https://shop.example.com/product/2",
        "https://shop.example.com/product/3",
    },
    Prompt:   "Extract product name, price, and availability",
    Output:   "json",
    UseProxy: true,
})
if err != nil {
    log.Fatal(err)
}

for _, item := range batch.Items {
    if item.Status == "completed" {
        fmt.Println(item.URL, item.Result)
    } else if item.Status == "failed" {
        fmt.Println(item.URL, "failed:", item.Error)
    }
}
```

**Retry failed items:**

```go
queued, err := client.Batch.Submit(ctx, spidra.BatchParams{
    URLs:   []string{"https://example.com/1", "https://example.com/2"},
    Prompt: "Extract the page title",
})
if err != nil {
    log.Fatal(err)
}

batchID := queued.BatchID

// Later, after checking status
result, err := client.Batch.Get(ctx, batchID)
if result.FailedCount > 0 {
    client.Batch.Retry(ctx, batchID)
}
```

**Cancel a running batch:**

```go
err := client.Batch.Cancel(ctx, batchID)
```

**List past batches:**

```go
result, err := client.Batch.List(ctx, 1, 20)
if err != nil {
    log.Fatal(err)
}

for _, job := range result.Jobs {
    fmt.Printf("%s %s %d/%d\n", job.UUID, job.Status, job.CompletedCount, job.TotalURLs)
}
```

## Crawling

Given a starting URL, Spidra discovers pages automatically according to your instruction and extracts structured data from each one.

```go
job, err := client.Crawl.Run(ctx, spidra.CrawlParams{
    BaseURL:              "https://competitor.com/blog",
    CrawlInstruction:     "Find all blog posts published in 2024",
    TransformInstruction: "Extract the title, author, publish date, and a one-sentence summary",
    MaxPages:             30,
    UseProxy:             true,
})
if err != nil {
    log.Fatal(err)
}

for _, page := range job.Result {
    fmt.Println(page.URL)
    fmt.Println(page.Data)
}
```

**Submit without waiting:**

```go
queued, err := client.Crawl.Submit(ctx, spidra.CrawlParams{
    BaseURL:              "https://example.com/docs",
    CrawlInstruction:     "Find all documentation pages",
    TransformInstruction: "Extract the page title and main content summary",
    MaxPages:             50,
})
if err != nil {
    log.Fatal(err)
}

jobID := queued.JobID

// Check status later
status, err := client.Crawl.Get(ctx, jobID)
```

**Get signed download URLs for all crawled pages:**

Each page includes `HTMLURL` and `MarkdownURL` pointing to S3-signed URLs that expire after 1 hour.

```go
result, err := client.Crawl.Pages(ctx, jobID)
if err != nil {
    log.Fatal(err)
}

for _, page := range result.Pages {
    fmt.Println(page.URL, "—", page.Status)
    // page.HTMLURL     — download raw HTML
    // page.MarkdownURL — download markdown
}
```

**Re-extract with a new instruction:**

Runs a new AI transformation over an existing completed crawl without re-crawling. Charges credits for the transformation only.

```go
newJob, err := client.Crawl.Extract(ctx, sourceJobID, "Extract only the product SKUs and prices as a flat list")
if err != nil {
    log.Fatal(err)
}

// Poll the new job
result, err := client.Crawl.Get(ctx, newJob.JobID)
```

**Crawl history and stats:**

```go
history, err := client.Crawl.History(ctx, 1, 10)
for _, job := range history.Jobs {
    fmt.Printf("%s — %s — %d pages\n", job.BaseURL, job.Status, job.PagesCrawled)
}

stats, err := client.Crawl.Stats(ctx)
fmt.Println("Total crawls:", stats.Total)
```

## Logs

Scrape logs are stored for every job that runs through the API.

```go
// List logs with optional filters
result, err := client.Logs.List(ctx, map[string]string{
    "status":     "failed",
    "searchTerm": "amazon.com",
    "dateStart":  "2024-01-01",
    "dateEnd":    "2024-12-31",
    "page":       "1",
    "limit":      "20",
})
if err != nil {
    log.Fatal(err)
}

for _, log := range result.Data.Logs {
    fmt.Println(log["status"], log["credits_used"])
}
```

**Get a single log with full extraction result:**

```go
log, err := client.Logs.Get(ctx, "log-uuid-here")
if err != nil {
    log.Fatal(err)
}
fmt.Println(log.Data["result_data"])
```

## Usage Statistics

Returns credit and request usage broken down by day or week.

```go
// Range options: "7d" | "30d" | "weekly"
result, err := client.Usage.Get(ctx, "30d")
if err != nil {
    log.Fatal(err)
}

for _, row := range result.Data {
    fmt.Printf("%s: %d requests, %.2f credits, %d tokens\n",
        row.Date, row.Requests, row.Credits, row.Tokens)
}
```

## Error Handling

Every API error is a typed value. Use `errors.As` to check for specific error types.

```go
import "errors"

job, err := client.Scrape.Run(ctx, params)
if err != nil {
    var authErr *spidra.AuthenticationError
    var credErr *spidra.InsufficientCreditsError
    var rateErr *spidra.RateLimitError
    var srvErr  *spidra.ServerError
    var apiErr  *spidra.SpidraError

    switch {
    case errors.As(err, &authErr):
        // 401 — API key is missing or invalid
        fmt.Println("Check your API key")
    case errors.As(err, &credErr):
        // 403 — monthly credit limit reached
        fmt.Println("Out of credits. Top up at app.spidra.io")
    case errors.As(err, &rateErr):
        // 429 — too many requests
        fmt.Println("Rate limited, back off and retry")
    case errors.As(err, &srvErr):
        // 5xx — something went wrong on Spidra's side
        fmt.Println("Server error, try again")
    case errors.As(err, &apiErr):
        // Any other API error
        fmt.Printf("%d: %s\n", apiErr.StatusCode, apiErr.Message)
    default:
        fmt.Println("Unexpected error:", err)
    }
}
```

## Custom Base URL

```go
client := spidra.New(
    "spd_YOUR_API_KEY",
    spidra.WithBaseURL("http://localhost:4321/api"), // for local development
)
```

## Resources

- [Spidra Documentation](https://docs.spidra.io)
- [Spidra Dashboard](https://app.spidra.io)
- [Report an issue](https://github.com/spidra-io/spidra-go/issues)

## License

MIT
