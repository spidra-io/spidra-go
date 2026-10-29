package spidra

import "context"

// LogsResource handles execution log retrieval across scrape, search, map, and parse.
type LogsResource struct {
	http *httpClient
}

// List returns scrape logs with optional filters (status, searchTerm, dateStart, dateEnd, page, limit).
func (r *LogsResource) List(ctx context.Context, filters map[string]string) (*LogsResult, error) {
	var result LogsResult
	if err := r.http.get(ctx, "/scrape-logs", filters, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Scrape is an alias for List.
func (r *LogsResource) Scrape(ctx context.Context, filters map[string]string) (*LogsResult, error) {
	return r.List(ctx, filters)
}

// Get returns a single scrape log with the full extraction result.
func (r *LogsResource) Get(ctx context.Context, logID string) (*LogResult, error) {
	var result LogResult
	if err := r.http.get(ctx, "/scrape-logs/"+logID, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetScrape is an alias for Get.
func (r *LogsResource) GetScrape(ctx context.Context, logID string) (*LogResult, error) {
	return r.Get(ctx, logID)
}

// Search returns search query execution logs with optional filters and pagination.
func (r *LogsResource) Search(ctx context.Context, queryParams map[string]string) (*SearchLogsResult, error) {
	var result SearchLogsResult
	if err := r.http.get(ctx, "/search-logs", queryParams, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSearch retrieves details of a single search log entry.
func (r *LogsResource) GetSearch(ctx context.Context, uuid string) (*SearchLog, error) {
	var result SearchLog
	if err := r.http.get(ctx, "/search-logs/"+uuid, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Map returns domain mapping execution logs with optional filters and pagination.
func (r *LogsResource) Map(ctx context.Context, queryParams map[string]string) (*MapLogsResult, error) {
	var result MapLogsResult
	if err := r.http.get(ctx, "/map-logs", queryParams, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMap retrieves details of a single domain map log entry.
func (r *LogsResource) GetMap(ctx context.Context, uuid string) (*MapLog, error) {
	var result MapLog
	if err := r.http.get(ctx, "/map-logs/"+uuid, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Parse returns document parsing execution logs with optional filters and pagination.
func (r *LogsResource) Parse(ctx context.Context, queryParams map[string]string) (*ParseLogsResult, error) {
	var result ParseLogsResult
	if err := r.http.get(ctx, "/parse-logs", queryParams, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
