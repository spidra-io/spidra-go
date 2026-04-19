package spidra

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// BatchResource handles batch scrape jobs (up to 50 URLs per request).
type BatchResource struct {
	http *httpClient
}

// List returns past batch scrape jobs, newest first.
func (r *BatchResource) List(ctx context.Context, page, limit int) (*BatchList, error) {
	var result BatchList
	err := r.http.get(ctx, "/batch/scrape", map[string]string{
		"page":  strconv.Itoa(page),
		"limit": strconv.Itoa(limit),
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Submit queues a batch scrape job and returns immediately with a BatchID.
func (r *BatchResource) Submit(ctx context.Context, params BatchParams) (*BatchJob, error) {
	var job BatchJob
	if err := r.http.post(ctx, "/batch/scrape", params, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// Get returns the current status and results of a batch scrape job.
func (r *BatchResource) Get(ctx context.Context, batchID string) (*BatchJob, error) {
	var job BatchJob
	if err := r.http.get(ctx, "/batch/scrape/"+batchID, nil, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// Retry re-queues all failed URLs in a batch scrape job.
func (r *BatchResource) Retry(ctx context.Context, batchID string) error {
	return r.http.post(ctx, "/batch/scrape/"+batchID+"/retry", nil, nil)
}

// Cancel stops a running batch scrape job and refunds credits for unprocessed URLs.
func (r *BatchResource) Cancel(ctx context.Context, batchID string) error {
	return r.http.delete(ctx, "/batch/scrape/"+batchID, nil)
}

// Run submits a batch scrape job and polls until it completes.
func (r *BatchResource) Run(ctx context.Context, params BatchParams, opts ...PollOptions) (*BatchJob, error) {
	o := PollOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	o = o.withDefaults()

	queued, err := r.Submit(ctx, params)
	if err != nil {
		return nil, err
	}

	batchID := queued.BatchID
	if batchID == "" {
		return nil, fmt.Errorf("spidra: no batchId in response")
	}

	deadline := time.Now().Add(o.Timeout)

	for {
		job, err := r.Get(ctx, batchID)
		if err != nil {
			return nil, err
		}

		switch job.Status {
		case "completed", "failed", "cancelled":
			job.BatchID = batchID
			return job, nil
		}

		if time.Now().Add(o.PollInterval).After(deadline) {
			return nil, fmt.Errorf("spidra: batch job timed out after %s, batchId: %s", o.Timeout, batchID)
		}

		time.Sleep(o.PollInterval)
	}
}
