package spidra

import "context"

// AccountResource provides access to account information and credit usage statistics.
type AccountResource struct {
	http *httpClient
}

// Get returns account subscription details, credit balances, and residential proxy limits.
func (r *AccountResource) Get(ctx context.Context) (*AccountSummary, error) {
	var summary AccountSummary
	if err := r.http.get(ctx, "/account", nil, &summary); err != nil {
		return nil, err
	}
	return &summary, nil
}

// Usage returns credit and request usage over time ("7d", "30d", "weekly").
func (r *AccountResource) Usage(ctx context.Context, dateRange string) ([]UsageStatRow, error) {
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
func (r *AccountResource) Breakdown(ctx context.Context, dateRange, startDate, endDate string) ([]UsageBreakdownRow, error) {
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

