package spidra

import "context"

// UsageResource handles credit and request usage statistics.
type UsageResource struct {
	http *httpClient
}

// Get returns usage statistics broken down by day or week ("7d", "30d", "weekly").
func (r *UsageResource) Get(ctx context.Context, dateRange string) (*UsageResult, error) {
	var result UsageResult
	err := r.http.get(ctx, "/usage-stats", map[string]string{"range": dateRange}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Stats returns structured credit and request usage records over time ("7d", "30d", "weekly").
func (r *UsageResource) Stats(ctx context.Context, dateRange string) ([]UsageStatRow, error) {
	if dateRange == "" {
		dateRange = "30d"
	}
	var rows []UsageStatRow
	if err := r.http.get(ctx, "/usage-stats", map[string]string{"range": dateRange}, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// Breakdown returns credit consumption categorized by tool ("1d", "7d", "30d", "custom").
func (r *UsageResource) Breakdown(ctx context.Context, dateRange, startDate, endDate string) ([]UsageBreakdownRow, error) {
	if dateRange == "" {
		dateRange = "7d"
	}
	query := map[string]string{"range": dateRange}
	if startDate != "" {
		query["startDate"] = startDate
	}
	if endDate != "" {
		query["endDate"] = endDate
	}
	var rows []UsageBreakdownRow
	if err := r.http.get(ctx, "/usage-stats/breakdown", query, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}
