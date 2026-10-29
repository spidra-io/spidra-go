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
	query := map[string]string{}
	if page > 0 {
		query["page"] = strconv.Itoa(page)
	}
	if limit > 0 {
		query["limit"] = strconv.Itoa(limit)
	}
	if err := r.http.get(ctx, "/batch/scrape", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Submit queues a batch scrape job and returns immediately with a BatchJob.
func (r *BatchResource) Submit(ctx context.Context, params BatchParams) (*BatchJob, error) {
	var job BatchJob
	if err := r.http.post(ctx, "/batch/scrape", params, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// SubmitURLs is a convenience helper to queue a list of URLs.
func (r *BatchResource) SubmitURLs(ctx context.Context, urls []string) (*BatchJob, error) {
	return r.Submit(ctx, BatchParams{URLs: urls})
}

// Get returns the current status and results of a batch scrape job.
func (r *BatchResource) Get(ctx context.Context, batchID string) (*BatchJob, error) {
	var job BatchJob
	if err := r.http.get(ctx, "/batch/scrape/"+batchID, nil, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// Retry re-queues all failed URLs in a batch scrape job with optional cookies.
func (r *BatchResource) Retry(ctx context.Context, batchID string, cookies ...string) error {
	var body any
	if len(cookies) > 0 && cookies[0] != "" {
		body = map[string]string{"cookies": cookies[0]}
	}
	return r.http.post(ctx, "/batch/scrape/"+batchID+"/retry", body, nil)
}

// Cancel stops a running batch scrape job and refunds credits for unprocessed URLs.
func (r *BatchResource) Cancel(ctx context.Context, batchID string) (*BatchCancelResult, error) {
	var result BatchCancelResult
	if err := r.http.delete(ctx, "/batch/scrape/"+batchID, &result); err != nil {
		return nil, err
	}
	return &result, nil
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

	switch queued.Status {
	case "completed":
		return queued, nil
	case "failed", "error", "cancelled":
		return queued, &JobFailedError{
			Message: fmt.Sprintf("batch job %s: failed %d of %d items", queued.Status, queued.FailedCount, queued.TotalURLs),
			JobID:   batchID,
			Status:  queued.Status,
		}
	}

	deadline := time.Now().Add(o.Timeout)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		job, err := r.Get(ctx, batchID)
		if err != nil {
			return nil, err
		}

		switch job.Status {
		case "completed":
			job.BatchID = batchID
			return job, nil
		case "failed", "error", "cancelled":
			job.BatchID = batchID
			return job, &JobFailedError{
				Message: fmt.Sprintf("batch job %s: failed %d of %d items", job.Status, job.FailedCount, job.TotalURLs),
				JobID:   batchID,
				Status:  job.Status,
			}
		}

		if time.Now().Add(o.PollInterval).After(deadline) {
			return nil, &TimeoutError{
				Message:        fmt.Sprintf("spidra: batch job timed out after %s, batchId: %s", o.Timeout, batchID),
				JobID:          batchID,
				TimeoutSeconds: int64(o.Timeout.Seconds()),
			}
		}

		time.Sleep(o.PollInterval)
	}
}

// RunURLs submits a slice of URLs and polls until all URLs finish scraping.
func (r *BatchResource) RunURLs(ctx context.Context, urls []string, opts ...PollOptions) (*BatchJob, error) {
	return r.Run(ctx, BatchParams{URLs: urls}, opts...)
}
