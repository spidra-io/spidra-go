package spidra

import "context"

// UsageResource handles credit and request usage statistics.
type UsageResource struct {
	http *httpClient
}

// Get returns usage statistics broken down by day or week.
// Range options: "7d", "30d", "weekly".
func (r *UsageResource) Get(ctx context.Context, dateRange string) (*UsageResult, error) {
	var result UsageResult
	err := r.http.get(ctx, "/usage-stats", map[string]string{"range": dateRange}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
