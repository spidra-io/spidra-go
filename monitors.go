package spidra

import "context"

// MonitorsResource handles scheduled website change monitors.
type MonitorsResource struct {
	http *httpClient
}

// List returns scheduled monitors with optional pagination parameters.
func (r *MonitorsResource) List(ctx context.Context, queryParams map[string]string) (*ListMonitorsResponse, error) {
	var resp ListMonitorsResponse
	if err := r.http.get(ctx, "/monitors", queryParams, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a scheduled monitor by its ID.
func (r *MonitorsResource) Get(ctx context.Context, monitorID string) (*Monitor, error) {
	var monitor Monitor
	if err := r.http.get(ctx, "/monitors/"+monitorID, nil, &monitor); err != nil {
		return nil, err
	}
	return &monitor, nil
}

// Create sets up a new scheduled website change monitor.
func (r *MonitorsResource) Create(ctx context.Context, params MonitorParams) (*Monitor, error) {
	var monitor Monitor
	if err := r.http.post(ctx, "/monitors", params, &monitor); err != nil {
		return nil, err
	}
	return &monitor, nil
}

// Update modifies an existing monitor configuration.
func (r *MonitorsResource) Update(ctx context.Context, monitorID string, updates map[string]any) (*Monitor, error) {
	var monitor Monitor
	if err := r.http.patch(ctx, "/monitors/"+monitorID, updates, &monitor); err != nil {
		return nil, err
	}
	return &monitor, nil
}

// Delete permanently removes a monitor by its ID.
func (r *MonitorsResource) Delete(ctx context.Context, monitorID string) error {
	return r.http.delete(ctx, "/monitors/"+monitorID, nil)
}

// Estimate calculates monthly credit cost projections for a monitor schedule.
func (r *MonitorsResource) Estimate(ctx context.Context, params map[string]any) (*MonitorEstimate, error) {
	var est MonitorEstimate
	if err := r.http.post(ctx, "/monitors/estimate", params, &est); err != nil {
		return nil, err
	}
	return &est, nil
}

// Limits returns active monitor limits and current usage for the account tier.
func (r *MonitorsResource) Limits(ctx context.Context) (*MonitorLimits, error) {
	var limits MonitorLimits
	if err := r.http.get(ctx, "/monitors/limits", nil, &limits); err != nil {
		return nil, err
	}
	return &limits, nil
}

// TestNotification sends a test alert across configured notification channels.
func (r *MonitorsResource) TestNotification(ctx context.Context, channels map[string]any, sampleData any) (*NotificationTestResult, error) {
	body := map[string]any{"channels": channels}
	if sampleData != nil {
		body["sampleData"] = sampleData
	}
	var result NotificationTestResult
	if err := r.http.post(ctx, "/monitors/test-notification", body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CheckNow triggers an immediate out-of-schedule check for a monitor.
func (r *MonitorsResource) CheckNow(ctx context.Context, monitorID string) (map[string]any, error) {
	var resp map[string]any
	if err := r.http.post(ctx, "/monitors/"+monitorID+"/check-now", map[string]any{}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Pause suspends scheduled checks for a monitor.
func (r *MonitorsResource) Pause(ctx context.Context, monitorID string) (*Monitor, error) {
	var monitor Monitor
	if err := r.http.post(ctx, "/monitors/"+monitorID+"/pause", map[string]any{}, &monitor); err != nil {
		return nil, err
	}
	return &monitor, nil
}

// Resume re-activates scheduled checks for a paused monitor.
func (r *MonitorsResource) Resume(ctx context.Context, monitorID string) (*Monitor, error) {
	var monitor Monitor
	if err := r.http.post(ctx, "/monitors/"+monitorID+"/resume", map[string]any{}, &monitor); err != nil {
		return nil, err
	}
	return &monitor, nil
}

// ListChecks returns historical check runs for a monitor with optional pagination.
func (r *MonitorsResource) ListChecks(ctx context.Context, monitorID string, queryParams map[string]string) (*ListMonitorChecksResponse, error) {
	var resp ListMonitorChecksResponse
	if err := r.http.get(ctx, "/monitors/"+monitorID+"/checks", queryParams, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCheck returns details of a single historical check run.
func (r *MonitorsResource) GetCheck(ctx context.Context, monitorID, checkID string) (*MonitorCheck, error) {
	var check MonitorCheck
	if err := r.http.get(ctx, "/monitors/"+monitorID+"/checks/"+checkID, nil, &check); err != nil {
		return nil, err
	}
	return &check, nil
}

