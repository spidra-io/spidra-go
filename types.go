package spidra

// ——— Scrape ———

// ScrapeParams are the parameters for a scrape job.
type ScrapeParams struct {
	URLs               []ScrapeURL     `json:"urls"`
	Prompt             string          `json:"prompt,omitempty"`
	Output             string          `json:"output,omitempty"`
	Schema             any             `json:"schema,omitempty"`
	UseProxy           bool            `json:"useProxy,omitempty"`
	ProxyCountry       string          `json:"proxyCountry,omitempty"`
	Screenshot         bool            `json:"screenshot,omitempty"`
	FullPageScreenshot bool            `json:"fullPageScreenshot,omitempty"`
	ExtractContentOnly bool            `json:"extractContentOnly,omitempty"`
	Cookies            string          `json:"cookies,omitempty"`
	Actions            []BrowserAction `json:"actions,omitempty"`
}

// BrowserAction specifies an interactive browser action (click, type, scroll, wait, forEach).
type BrowserAction struct {
	Type         string          `json:"type"`
	Selector     string          `json:"selector,omitempty"`
	Text         string          `json:"text,omitempty"`
	Milliseconds int             `json:"milliseconds,omitempty"`
	Actions      []BrowserAction `json:"actions,omitempty"`
	Parameters   map[string]any  `json:"parameters,omitempty"`
}

// ScrapeURL is a single URL entry with optional browser actions.
type ScrapeURL struct {
	URL     string           `json:"url"`
	Actions []map[string]any `json:"actions,omitempty"`
}

// ScrapeJob is the response from a scrape status check or completed run.
type ScrapeJob struct {
	JobID  string        `json:"jobId"`
	Status string        `json:"status"`
	Result *ScrapeOutput `json:"result,omitempty"`
	Error  string        `json:"error,omitempty"`
}

// ScrapeOutput holds the extraction result of a completed scrape job.
type ScrapeOutput struct {
	Content     any            `json:"content,omitempty"`
	Screenshots []string       `json:"screenshots,omitempty"`
	Actions     []ActionResult `json:"actions,omitempty"`
}

// ActionResult holds the execution result of a single browser action.
type ActionResult struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	DurationMs int64  `json:"durationMs,omitempty"`
	Error      string `json:"error,omitempty"`
}

// ——— Search ———

// SearchParams are the parameters for a web/news/multimodal search job.
type SearchParams struct {
	Query          string              `json:"query"`
	Sources        []string            `json:"sources,omitempty"` // web, news, images, videos, research, developer
	Limit          int                 `json:"limit,omitempty"`
	Domain         string              `json:"domain,omitempty"`
	Domains        []string            `json:"domains,omitempty"`
	ExcludeDomains []string            `json:"excludeDomains,omitempty"`
	DateStart      string              `json:"dateStart,omitempty"`
	DateEnd        string              `json:"dateEnd,omitempty"`
	Recency        string              `json:"recency,omitempty"`
	Country        string              `json:"country,omitempty"`
	Language       string              `json:"language,omitempty"`
	Scrape         bool                `json:"scrape,omitempty"`
	ScrapeOptions  *SearchScrapeOptions `json:"scrapeOptions,omitempty"`
	Wait           bool                `json:"wait,omitempty"`
}

// SearchScrapeOptions controls opt-in markdown scraping for search results.
type SearchScrapeOptions struct {
	Output       string `json:"output,omitempty"`
	Prompt       string `json:"prompt,omitempty"`
	Schema       any    `json:"schema,omitempty"`
	UseProxy     bool   `json:"useProxy,omitempty"`
	ProxyCountry string `json:"proxyCountry,omitempty"`
}

// SearchJob is the response from a search status check or completed query.
type SearchJob struct {
	JobID  string        `json:"jobId"`
	Status string        `json:"status"`
	Query  string        `json:"query,omitempty"`
	Result *SearchResult `json:"result,omitempty"`
	Error  string        `json:"error,omitempty"`
}

// SearchResult holds organized search items, extras, and stats.
type SearchResult struct {
	Query   string              `json:"query,omitempty"`
	Data    *SearchResponseData `json:"data,omitempty"`
	Results []SearchResultItem  `json:"results,omitempty"`
	Extras  *SearchExtras       `json:"extras,omitempty"`
	Stats   *SearchStats        `json:"stats,omitempty"`
}

// SearchResponseData groups search results by source category.
type SearchResponseData struct {
	Web       []SearchResultItem `json:"web,omitempty"`
	News      []SearchResultItem `json:"news,omitempty"`
	Images    []SearchResultItem `json:"images,omitempty"`
	Videos    []SearchResultItem `json:"videos,omitempty"`
	Research  []SearchResultItem `json:"research,omitempty"`
	Developer []SearchResultItem `json:"developer,omitempty"`
}

// SearchResultItem represents an individual search result.
type SearchResultItem struct {
	Title          string `json:"title,omitempty"`
	URL            string `json:"url,omitempty"`
	Snippet        string `json:"snippet,omitempty"`
	Date           string `json:"date,omitempty"`
	Source         string `json:"source,omitempty"`
	Image          string `json:"image,omitempty"`
	Thumbnail      string `json:"thumbnail,omitempty"`
	ScrapedContent any    `json:"scrapedContent,omitempty"`
}

// SearchExtras contains AI synthesized answers and related questions.
type SearchExtras struct {
	Answer           string   `json:"answer,omitempty"`
	RelatedQuestions []string `json:"relatedQuestions,omitempty"`
}

// SearchStats contains execution metrics for the search query.
type SearchStats struct {
	TotalResults    int   `json:"totalResults,omitempty"`
	ExecutionTimeMs int64 `json:"executionTimeMs,omitempty"`
}

// ——— Map ———

// MapParams are the parameters for a domain mapping job.
type MapParams struct {
	URL               string `json:"url"`
	Search            string `json:"search,omitempty"`
	IgnoreSitemap     bool   `json:"ignoreSitemap,omitempty"`
	IncludeSubdomains bool   `json:"includeSubdomains,omitempty"`
	Limit             int    `json:"limit,omitempty"`
	UseProxy          bool   `json:"useProxy,omitempty"`
	ProxyCountry      string `json:"proxyCountry,omitempty"`
}

// MapResult is the response from a domain mapping request.
type MapResult struct {
	Links           []MapLink `json:"links,omitempty"`
	TotalLinks      int       `json:"totalLinks,omitempty"`
	ExecutionTimeMs int64     `json:"executionTimeMs,omitempty"`
}

// MapLink is a discovered page link on the target domain.
type MapLink struct {
	URL         string `json:"url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// ——— Parse ———

// ParseParams are the parameters for a document or media parsing request.
type ParseParams struct {
	URL          string `json:"url,omitempty"`
	FileID       string `json:"fileId,omitempty"`
	FileName     string `json:"fileName,omitempty"`
	Pages        string `json:"pages,omitempty"`
	IncludePages []int  `json:"includePages,omitempty"`
	Prompt       string `json:"prompt,omitempty"`
	Mode         string `json:"mode,omitempty"` // "ai", "fast"
	Wait         bool   `json:"wait,omitempty"`
}

// ParseJob is the response from a document parsing job.
type ParseJob struct {
	JobID           string         `json:"jobId"`
	Status          string         `json:"status"`
	Result          *ParsedContent `json:"result,omitempty"`
	Pages           []ParsedPage   `json:"pages,omitempty"`
	PageCount       int            `json:"pageCount,omitempty"`
	TotalCredits    float64        `json:"totalCredits,omitempty"`
	ExecutionTimeMs int64          `json:"executionTimeMs,omitempty"`
	Error           string         `json:"error,omitempty"`
}

// ParsedContent holds aggregate parsed text, markdown, or structured data.
type ParsedContent struct {
	Markdown string `json:"markdown,omitempty"`
	Text     string `json:"text,omitempty"`
	Data     any    `json:"data,omitempty"`
}

// ParsedPage holds page-level parsed output and diagrams/images.
type ParsedPage struct {
	PageNumber int      `json:"pageNumber,omitempty"`
	Markdown   string   `json:"markdown,omitempty"`
	Text       string   `json:"text,omitempty"`
	Images     []string `json:"images,omitempty"`
}

// ParseEstimate provides credit cost and timing estimates for document parsing.
type ParseEstimate struct {
	EstimatedCredits   float64 `json:"estimatedCredits,omitempty"`
	PageCount          int     `json:"pageCount,omitempty"`
	DurationEstimateMs int64   `json:"durationEstimateMs,omitempty"`
}

// ParseBatch represents a multi-document parsing batch.
type ParseBatch struct {
	BatchID        string     `json:"batchId"`
	Status         string     `json:"status"`
	TotalFiles     int        `json:"totalFiles,omitempty"`
	CompletedCount int        `json:"completedCount,omitempty"`
	FailedCount    int        `json:"failedCount,omitempty"`
	Items          []ParseJob `json:"items,omitempty"`
}

// ——— Turnkey Templates ———

// TemplateSummary provides an overview of an available turnkey extraction template.
type TemplateSummary struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
	Site        string `json:"site,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

// TemplateDetail contains the full schema, field inputs, and pricing for a turnkey template.
type TemplateDetail struct {
	Slug              string           `json:"slug"`
	Name              string           `json:"name"`
	Description       string           `json:"description,omitempty"`
	Category          string           `json:"category,omitempty"`
	Site              string           `json:"site,omitempty"`
	Inputs            []TemplateInput  `json:"inputs,omitempty"`
	Columns           []TemplateColumn `json:"columns,omitempty"`
	BaseCreditCost    float64          `json:"baseCreditCost,omitempty"`
	CreditCostFormula string           `json:"creditCostFormula,omitempty"`
}

// TemplateInput describes a parameter field required by a turnkey template.
type TemplateInput struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Label       string `json:"label,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Default     any    `json:"default,omitempty"`
	Options     []any  `json:"options,omitempty"`
	Description string `json:"description,omitempty"`
}

// TemplateColumn describes a data column extracted by a turnkey template.
type TemplateColumn struct {
	Key   string `json:"key"`
	Label string `json:"label,omitempty"`
	Type  string `json:"type,omitempty"`
}

// TemplateSite describes a supported website/platform for turnkey templates.
type TemplateSite struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	Icon string `json:"icon,omitempty"`
}

// TemplateEstimate provides the credit cost estimate for running a turnkey template.
type TemplateEstimate struct {
	EstimatedCredits float64        `json:"estimatedCredits"`
	Breakdown        map[string]any `json:"breakdown,omitempty"`
}

// TemplateRun is the status and result of a turnkey template execution.
type TemplateRun struct {
	RunID       string               `json:"runId"`
	Slug        string               `json:"slug,omitempty"`
	Status      string               `json:"status"`
	Inputs      map[string]any       `json:"inputs,omitempty"`
	Result      []map[string]any     `json:"result,omitempty"`
	Progress    *TemplateRunProgress `json:"progress,omitempty"`
	CreditsUsed float64              `json:"creditsUsed,omitempty"`
	TotalRows   int                  `json:"totalRows,omitempty"`
	Error       string               `json:"error,omitempty"`
	CreatedAt   string               `json:"createdAt,omitempty"`
	CompletedAt string               `json:"completedAt,omitempty"`
}

// TemplateRunProgress indicates running template progress.
type TemplateRunProgress struct {
	Total     int     `json:"total,omitempty"`
	Completed int     `json:"completed,omitempty"`
	Failed    int     `json:"failed,omitempty"`
	Percent   float64 `json:"percent,omitempty"`
}

// Airport represents airport data from the airport lookup directory.
type Airport struct {
	IATA    string `json:"iata"`
	ICAO    string `json:"icao,omitempty"`
	Name    string `json:"name"`
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
}

// ListTemplatesResponse is the response when listing available turnkey templates.
type ListTemplatesResponse struct {
	Templates  []TemplateSummary `json:"templates"`
	Categories []string          `json:"categories,omitempty"`
	Sites      []TemplateSite    `json:"sites,omitempty"`
	Total      int               `json:"total,omitempty"`
}

// ListTemplateRunsResponse is the response when listing past template executions.
type ListTemplateRunsResponse struct {
	Runs  []TemplateRun `json:"runs"`
	Total int           `json:"total,omitempty"`
}

// ——— Presets ———

// Preset represents a saved reusable tool or template configuration.
type Preset struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Tool         string         `json:"tool"`
	Description  string         `json:"description,omitempty"`
	TemplateSlug string         `json:"templateSlug,omitempty"`
	Config       map[string]any `json:"config,omitempty"`
	CreatedAt    string         `json:"createdAt,omitempty"`
	UpdatedAt    string         `json:"updatedAt,omitempty"`
}

// PresetParams are the parameters to create or save a preset.
type PresetParams struct {
	Name         string         `json:"name"`
	Tool         string         `json:"tool"`
	Description  string         `json:"description,omitempty"`
	TemplateSlug string         `json:"templateSlug,omitempty"`
	Config       map[string]any `json:"config,omitempty"`
}

// RunPresetResponse represents the result of executing a saved preset.
type RunPresetResponse struct {
	Success bool   `json:"success"`
	JobID   string `json:"jobId,omitempty"`
	BatchID string `json:"batchId,omitempty"`
	RunID   string `json:"runId,omitempty"`
	Status  string `json:"status,omitempty"`
	Result  any    `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ListPresetsResponse is the response from listing saved presets.
type ListPresetsResponse struct {
	Presets []Preset `json:"presets"`
	Total   int      `json:"total,omitempty"`
}

// ——— Monitors ———

// Monitor represents a scheduled website change monitor.
type Monitor struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	URL           string         `json:"url"`
	Schedule      string         `json:"schedule"`
	Status        string         `json:"status"` // active, paused
	Selector      string         `json:"selector,omitempty"`
	Condition     string         `json:"condition,omitempty"`
	Channels      map[string]any `json:"channels,omitempty"`
	LastCheckedAt string         `json:"lastCheckedAt,omitempty"`
	LastChangedAt string         `json:"lastChangedAt,omitempty"`
	CreatedAt     string         `json:"createdAt,omitempty"`
	UpdatedAt     string         `json:"updatedAt,omitempty"`
}

// MonitorParams are the parameters to create or update a monitor.
type MonitorParams struct {
	Name         string         `json:"name"`
	URL          string         `json:"url"`
	Schedule     string         `json:"schedule"`
	Selector     string         `json:"selector,omitempty"`
	Condition    string         `json:"condition,omitempty"`
	Channels     map[string]any `json:"channels,omitempty"`
	UseProxy     bool           `json:"useProxy,omitempty"`
	ProxyCountry string         `json:"proxyCountry,omitempty"`
}

// MonitorEstimate provides credit cost estimations for a scheduled monitor.
type MonitorEstimate struct {
	EstimatedMonthlyCredits float64 `json:"estimatedMonthlyCredits"`
	ChecksPerMonth          int     `json:"checksPerMonth"`
	CostPerCheck            float64 `json:"costPerCheck"`
}

// MonitorLimits details account limits and active monitor usage.
type MonitorLimits struct {
	ActiveMonitors         int `json:"activeMonitors"`
	MaxMonitors            int `json:"maxMonitors"`
	MinimumIntervalMinutes int `json:"minimumIntervalMinutes"`
}

// NotificationTestResult details delivery status across configured notification channels.
type NotificationTestResult struct {
	Success  bool              `json:"success"`
	Channels map[string]bool   `json:"channels,omitempty"`
	Errors   map[string]string `json:"errors,omitempty"`
}

// MonitorCheck represents a single check execution for a monitor.
type MonitorCheck struct {
	ID         string                `json:"id"`
	MonitorID  string                `json:"monitorId"`
	CheckedAt  string                `json:"checkedAt"`
	HasChanges bool                  `json:"hasChanges"`
	Diff       string                `json:"diff,omitempty"`
	Changes    *MonitorSourceChanges `json:"changes,omitempty"`
	Error      string                `json:"error,omitempty"`
}

// MonitorSourceChanges describes the raw and structured changes detected during a check.
type MonitorSourceChanges struct {
	PreviousContent  string `json:"previousContent,omitempty"`
	CurrentContent   string `json:"currentContent,omitempty"`
	ExtractedChanges any    `json:"extractedChanges,omitempty"`
}

// ListMonitorsResponse is the response from listing scheduled monitors.
type ListMonitorsResponse struct {
	Monitors []Monitor `json:"monitors"`
	Total    int       `json:"total,omitempty"`
}

// ListMonitorChecksResponse is the response from listing historical checks for a monitor.
type ListMonitorChecksResponse struct {
	Checks []MonitorCheck `json:"checks"`
	Total  int            `json:"total,omitempty"`
}

// ——— Account & Usage ———

// AccountSummary holds user plan, credit balances, and residential proxy limits.
type AccountSummary struct {
	Email          string                 `json:"email"`
	Plan           *AccountPlan           `json:"plan,omitempty"`
	Credits        *AccountCredits        `json:"credits,omitempty"`
	ProxyBandwidth *AccountProxyBandwidth `json:"proxyBandwidth,omitempty"`
}

// AccountPlan details the active subscription tier.
type AccountPlan struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	RenewsAt string `json:"renewsAt,omitempty"`
}

// AccountCredits details available and total credits.
type AccountCredits struct {
	Available float64 `json:"available"`
	Total     float64 `json:"total"`
	Used      float64 `json:"used"`
}

// AccountProxyBandwidth details available and consumed residential proxy bandwidth in MB.
type AccountProxyBandwidth struct {
	AvailableMB float64 `json:"availableMb"`
	TotalMB     float64 `json:"totalMb"`
	UsedMB      float64 `json:"usedMb"`
}

// UsageStatRow holds credit and request usage for a single date period.
type UsageStatRow struct {
	Date     string  `json:"date"`
	Requests int     `json:"requests"`
	Credits  float64 `json:"credits"`
	Tokens   int     `json:"tokens"`
}

// UsageBreakdownRow holds credit consumption categorized by tool.
type UsageBreakdownRow struct {
	Date     string  `json:"date"`
	Tool     string  `json:"tool"`
	Credits  float64 `json:"credits"`
	Requests int     `json:"requests"`
}

// UsageResult is the legacy response from the usage statistics endpoint.
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

// BatchCancelResult is the response from cancelling a batch scrape job.
type BatchCancelResult struct {
	Success         bool    `json:"success"`
	RefundedCredits float64 `json:"refundedCredits,omitempty"`
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

// CrawlCancelResult is the response when cancelling an active crawl job.
type CrawlCancelResult struct {
	Success         bool    `json:"success"`
	Message         string  `json:"message,omitempty"`
	RefundedCredits float64 `json:"refundedCredits,omitempty"`
}

// CrawlReport holds skipped URLs and crawl error diagnostics.
type CrawlReport struct {
	SkippedURLs []map[string]any `json:"skippedUrls,omitempty"`
	Errors      []map[string]any `json:"errors,omitempty"`
}

// CrawlReportResponse is the response from the crawl diagnostics report endpoint.
type CrawlReportResponse struct {
	JobID  string       `json:"jobId"`
	Report *CrawlReport `json:"report,omitempty"`
}

// CrawlRetryPageResponse is the response from retrying a failed crawl page.
type CrawlRetryPageResponse struct {
	Success bool   `json:"success"`
	PageID  string `json:"pageId"`
	Status  string `json:"status"`
}

// ——— Logs ———

// LogsResult is the response from listing scrape logs.
type LogsResult struct {
	Status string   `json:"status"`
	Data   LogsData `json:"data"`
}

// LogsData holds scrape log entries and total count.
type LogsData struct {
	Logs  []map[string]any `json:"logs"`
	Total int              `json:"total"`
}

// LogResult is the response from fetching a single scrape log.
type LogResult struct {
	Status string         `json:"status"`
	Data   map[string]any `json:"data"`
}

// SearchLog represents a single search log entry.
type SearchLog struct {
	UUID      string   `json:"uuid"`
	Query     string   `json:"query"`
	Sources   []string `json:"sources,omitempty"`
	CreatedAt string   `json:"createdAt,omitempty"`
	Credits   float64  `json:"credits,omitempty"`
}

// SearchLogsResult is the response when listing search logs.
type SearchLogsResult struct {
	Status string `json:"status"`
	Data   struct {
		Logs  []SearchLog `json:"logs"`
		Total int         `json:"total"`
	} `json:"data"`
}

// MapLog represents a single domain map log entry.
type MapLog struct {
	UUID       string  `json:"uuid"`
	URL        string  `json:"url"`
	TotalLinks int     `json:"totalLinks,omitempty"`
	CreatedAt  string  `json:"createdAt,omitempty"`
	Credits    float64 `json:"credits,omitempty"`
}

// MapLogsResult is the response when listing domain map logs.
type MapLogsResult struct {
	Status string `json:"status"`
	Data   struct {
		Logs  []MapLog `json:"logs"`
		Total int      `json:"total"`
	} `json:"data"`
}

// ParseLog represents a single document parse log entry.
type ParseLog struct {
	UUID      string  `json:"uuid"`
	URL       string  `json:"url,omitempty"`
	FileName  string  `json:"fileName,omitempty"`
	Mode      string  `json:"mode,omitempty"`
	PageCount int     `json:"pageCount,omitempty"`
	CreatedAt string  `json:"createdAt,omitempty"`
	Credits   float64 `json:"credits,omitempty"`
}

// ParseLogsResult is the response when listing document parse logs.
type ParseLogsResult struct {
	Status string `json:"status"`
	Data   struct {
		Logs  []ParseLog `json:"logs"`
		Total int        `json:"total"`
	} `json:"data"`
}
