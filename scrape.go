package spidra

import (
	"context"
	"fmt"
	"time"
)

// ScrapeResource handles single and multi-URL scrape jobs.
type ScrapeResource struct {
	http *httpClient
}

// Submit queues a scrape job and returns immediately with a JobID.
func (r *ScrapeResource) Submit(ctx context.Context, params ScrapeParams) (*ScrapeJob, error) {
	var job ScrapeJob
	if err := r.http.post(ctx, "/scrape", params, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// Get returns the current status and result of a scrape job.
func (r *ScrapeResource) Get(ctx context.Context, jobID string) (*ScrapeJob, error) {
	var job ScrapeJob
	if err := r.http.get(ctx, "/scrape/"+jobID, nil, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// Run submits a scrape job and polls until it completes.
func (r *ScrapeResource) Run(ctx context.Context, params ScrapeParams, opts ...PollOptions) (*ScrapeJob, error) {
	o := PollOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	o = o.withDefaults()

	queued, err := r.Submit(ctx, params)
	if err != nil {
		return nil, err
	}

	jobID := queued.JobID
	if jobID == "" {
		return nil, fmt.Errorf("spidra: no jobId in response")
	}

	deadline := time.Now().Add(o.Timeout)

	for {
		job, err := r.Get(ctx, jobID)
		if err != nil {
			return nil, err
		}

		switch job.Status {
		case "completed", "failed", "cancelled":
			job.JobID = jobID
			return job, nil
		}

		if time.Now().Add(o.PollInterval).After(deadline) {
			return nil, fmt.Errorf("spidra: scrape job timed out after %s, jobId: %s", o.Timeout, jobID)
		}

		time.Sleep(o.PollInterval)
	}
}
