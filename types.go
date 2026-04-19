package spidra

// ——— Scrape ———

// ScrapeParams are the parameters for a scrape job.
type ScrapeParams struct {
	URLs               []ScrapeURL `json:"urls"`
	Prompt             string      `json:"prompt,omitempty"`
	Output             string      `json:"output,omitempty"`
	Schema             any         `json:"schema,omitempty"`
	UseProxy           bool        `json:"useProxy,omitempty"`
	ProxyCountry       string      `json:"proxyCountry,omitempty"`
	Screenshot         bool        `json:"screenshot,omitempty"`
	FullPageScreenshot bool        `json:"fullPageScreenshot,omitempty"`
	ExtractContentOnly bool        `json:"extractContentOnly,omitempty"`
	Cookies            string      `json:"cookies,omitempty"`
}

// ScrapeURL is a single URL entry with optional browser actions.
// Actions can be click, type, scroll, wait, or forEach — see the Spidra docs.
type ScrapeURL struct {
	URL     string           `json:"url"`
	Actions []map[string]any `json:"actions,omitempty"`
}

// ScrapeJob is the response from a scrape status check or completed run.
type ScrapeJob struct {
	JobID  string        `json:"jobId"`
	Status string        `json:"status"`
	Result *ScrapeOutput `json:"result"`
	Error  string        `json:"error"`
}

// ScrapeOutput holds the extraction result of a completed scrape job.
// Content is a string for markdown output, an object or array for JSON output.
type ScrapeOutput struct {
	Content     any      `json:"content"`
	Screenshots []string `json:"screenshots"`
}

// ——— Batch ———

// BatchParams are the parameters for a batch scrape job.
type BatchParams struct {
	URLs               []string `json:"urls"`
	Prompt             string   `json:"prompt,omitempty"`
	Output             string   `json:"output,omitempty"`
	Schema             any      `json:"schema,omitempty"`
	UseProxy           bool     `json:"useProxy,omitempty"`
	ProxyCountry       string   `json:"proxyCountry,omitempty"`
	Screenshot         bool     `json:"screenshot,omitempty"`
	FullPageScreenshot bool     `json:"fullPageScreenshot,omitempty"`
	ExtractContentOnly bool     `json:"extractContentOnly,omitempty"`
	Cookies            string   `json:"cookies,omitempty"`
}

// BatchJob is the response from a batch scrape status check or completed run.
type BatchJob struct {
	BatchID        string      `json:"batchId"`
	Status         string      `json:"status"`
	TotalURLs      int         `json:"totalUrls"`
	CompletedCount int         `json:"completedCount"`
	FailedCount    int         `json:"failedCount"`
	Items          []BatchItem `json:"items"`
}

// BatchItem is the result for a single URL within a batch job.
type BatchItem struct {
	URL    string `json:"url"`
	Status string `json:"status"`
	Result any    `json:"result"`
	Error  string `json:"error"`
}

// BatchList is the response from listing past batch jobs.
type BatchList struct {
	Jobs []BatchSummary `json:"jobs"`
}

// BatchSummary is a brief overview of a batch job returned in list responses.
type BatchSummary struct {
	UUID           string `json:"uuid"`
	Status         string `json:"status"`
	TotalURLs      int    `json:"totalUrls"`
	CompletedCount int    `json:"completedCount"`
	FailedCount    int    `json:"failedCount"`
}

// ——— Crawl ———

// CrawlParams are the parameters for a crawl job.
type CrawlParams struct {
	BaseURL              string `json:"baseUrl"`
	CrawlInstruction     string `json:"crawlInstruction"`
	TransformInstruction string `json:"transformInstruction"`
	MaxPages             int    `json:"maxPages,omitempty"`
	UseProxy             bool   `json:"useProxy,omitempty"`
	ProxyCountry         string `json:"proxyCountry,omitempty"`
	Cookies              string `json:"cookies,omitempty"`
}

// CrawlJob is the response from a crawl status check or completed run.
type CrawlJob struct {
	JobID    string      `json:"jobId"`
	Status   string      `json:"status"`
	Progress any         `json:"progress"`
	Result   []CrawlPage `json:"result"`
	Error    string      `json:"error"`
}

// CrawlPage is a single page collected during a crawl.
// Data holds the AI-extracted content and can be a string or structured object.
type CrawlPage struct {
	URL         string `json:"url"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	Data        any    `json:"data"`
	HTMLURL     string `json:"html_url"`
	MarkdownURL string `json:"markdown_url"`
	HTMLKey     string `json:"htmlKey"`
	MarkdownKey string `json:"markdownKey"`
}

// CrawlPages is the response from fetching all pages of a completed crawl.
type CrawlPages struct {
	Pages []CrawlPage `json:"pages"`
}

// CrawlHistory is the response from listing past crawl jobs.
type CrawlHistory struct {
	Jobs []CrawlSummary `json:"jobs"`
}

// CrawlSummary is a brief overview of a crawl job returned in history responses.
type CrawlSummary struct {
	UUID         string `json:"uuid"`
	BaseURL      string `json:"base_url"`
	Status       string `json:"status"`
	PagesCrawled int    `json:"pages_crawled"`
}

// CrawlStats holds the total crawl job count for the account.
type CrawlStats struct {
	Total int `json:"total"`
}

// CrawlExtract is the response when a new extraction job is queued.
type CrawlExtract struct {
	JobID string `json:"jobId"`
}

// ——— Logs ———

// LogsResult is the response from listing scrape logs.
type LogsResult struct {
	Status string   `json:"status"`
	Data   LogsData `json:"data"`
}

// LogsData holds the log entries and total count.
type LogsData struct {
	Logs  []map[string]any `json:"logs"`
	Total int              `json:"total"`
}

// LogResult is the response from fetching a single scrape log.
type LogResult struct {
	Status string         `json:"status"`
	Data   map[string]any `json:"data"`
}

// ——— Usage ———

// UsageResult is the response from the usage statistics endpoint.
type UsageResult struct {
	Data []UsageRow `json:"data"`
}

// UsageRow holds credit and request usage for a single day or week.
type UsageRow struct {
	Date     string  `json:"date"`
	Requests int     `json:"requests"`
	Credits  float64 `json:"credits"`
	Tokens   int     `json:"tokens"`
}
