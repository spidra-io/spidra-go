package spidra

import "context"

// LogsResource handles scrape log retrieval.
type LogsResource struct {
	http *httpClient
}

// List returns scrape logs with optional filters.
// Supported filter keys: status, searchTerm, dateStart, dateEnd, page, limit.
func (r *LogsResource) List(ctx context.Context, filters map[string]string) (*LogsResult, error) {
	var result LogsResult
	if err := r.http.get(ctx, "/scrape-logs", filters, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a single scrape log with the full extraction result.
func (r *LogsResource) Get(ctx context.Context, logID string) (*LogResult, error) {
	var result LogResult
	if err := r.http.get(ctx, "/scrape-logs/"+logID, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
