package spidra

import (
	"context"
	"fmt"
	"time"
)

// TemplatesResource handles Spidra's 100+ turnkey extraction templates.
type TemplatesResource struct {
	http *httpClient
}

// List returns available turnkey templates with optional filtering (category, site, search).
func (r *TemplatesResource) List(ctx context.Context, filters map[string]string) (*ListTemplatesResponse, error) {
	var resp ListTemplatesResponse
	if err := r.http.get(ctx, "/templates", filters, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns the full detail, input schema, and columns for a template.
func (r *TemplatesResource) Get(ctx context.Context, slug string) (*TemplateDetail, error) {
	var detail TemplateDetail
	if err := r.http.get(ctx, "/templates/"+slug, nil, &detail); err != nil {
		return nil, err
	}
	return &detail, nil
}

// Estimate calculates the estimated credit cost for running a template with the given inputs.
func (r *TemplatesResource) Estimate(ctx context.Context, slug string, inputs map[string]any) (*TemplateEstimate, error) {
	if inputs == nil {
		inputs = map[string]any{}
	}
	body := map[string]any{"inputs": inputs}
	var estimate TemplateEstimate
	if err := r.http.post(ctx, "/templates/"+slug+"/estimate", body, &estimate); err != nil {
		return nil, err
	}
	return &estimate, nil
}

// SubmitRun queues a template execution job and returns immediately with a TemplateRun.
func (r *TemplatesResource) SubmitRun(ctx context.Context, slug string, inputs map[string]any) (*TemplateRun, error) {
	if inputs == nil {
		inputs = map[string]any{}
	}
	body := map[string]any{"inputs": inputs}
	var run TemplateRun
	if err := r.http.post(ctx, "/templates/"+slug+"/run", body, &run); err != nil {
		return nil, err
	}
	return &run, nil
}

// GetRun retrieves the status and extracted data of a specific template run.
func (r *TemplatesResource) GetRun(ctx context.Context, slug, runID string) (*TemplateRun, error) {
	var run TemplateRun
	if err := r.http.get(ctx, "/templates/"+slug+"/runs/"+runID, nil, &run); err != nil {
		return nil, err
	}
	return &run, nil
}

// Run executes a turnkey template and waits for completion.
func (r *TemplatesResource) Run(ctx context.Context, slug string, inputs map[string]any, opts ...PollOptions) (*TemplateRun, error) {
	o := PollOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	o = o.withDefaults()

	if inputs == nil {
		inputs = map[string]any{}
	}
	body := map[string]any{
		"inputs": inputs,
		"wait":   true,
	}

	var run TemplateRun
	if err := r.http.post(ctx, "/templates/"+slug+"/run", body, &run); err != nil {
		return nil, err
	}

	runID := run.RunID

	switch run.Status {
	case "completed":
		return &run, nil
	case "failed", "cancelled":
		errMsg := run.Error
		if errMsg == "" {
			errMsg = "template run " + run.Status
		}
		return &run, &JobFailedError{Message: errMsg, JobID: runID, Status: run.Status}
	}

	if runID == "" {
		return &run, nil
	}

	deadline := time.Now().Add(o.Timeout)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		current, err := r.GetRun(ctx, slug, runID)
		if err != nil {
			return nil, err
		}

		switch current.Status {
		case "completed":
			current.RunID = runID
			return current, nil
		case "failed", "cancelled":
			current.RunID = runID
			errMsg := current.Error
			if errMsg == "" {
				errMsg = "template run " + current.Status
			}
			return current, &JobFailedError{Message: errMsg, JobID: runID, Status: current.Status}
		}

		if time.Now().Add(o.PollInterval).After(deadline) {
			return nil, &TimeoutError{
				Message:        fmt.Sprintf("spidra: template run timed out after %s, runId: %s", o.Timeout, runID),
				JobID:          runID,
				TimeoutSeconds: int64(o.Timeout.Seconds()),
			}
		}

		time.Sleep(o.PollInterval)
	}
}

// ListRuns returns past executions for a given template.
func (r *TemplatesResource) ListRuns(ctx context.Context, slug string, queryParams map[string]string) (*ListTemplateRunsResponse, error) {
	var resp ListTemplateRunsResponse
	if err := r.http.get(ctx, "/templates/"+slug+"/runs", queryParams, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupAirports searches airports by name, city, or IATA code.
func (r *TemplatesResource) LookupAirports(ctx context.Context, query string) ([]Airport, error) {
	var airports []Airport
	params := map[string]string{"q": query}
	if err := r.http.get(ctx, "/templates/lookup/airports", params, &airports); err != nil {
		return nil, err
	}
	return airports, nil
}

