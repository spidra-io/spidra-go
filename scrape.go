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

// SubmitURL is a convenience helper to queue a single URL with default options.
func (r *ScrapeResource) SubmitURL(ctx context.Context, targetURL string) (*ScrapeJob, error) {
	return r.Submit(ctx, ScrapeParams{
		URLs: []ScrapeURL{{URL: targetURL}},
	})
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

	switch queued.Status {
	case "completed":
		return queued, nil
	case "failed", "error", "cancelled":
		errMsg := queued.Error
		if errMsg == "" {
			errMsg = "scrape job " + queued.Status
		}
		return queued, &JobFailedError{Message: errMsg, JobID: jobID, Status: queued.Status}
	}

	deadline := time.Now().Add(o.Timeout)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		job, err := r.Get(ctx, jobID)
		if err != nil {
			return nil, err
		}

		switch job.Status {
		case "completed":
			job.JobID = jobID
			return job, nil
		case "failed", "error", "cancelled":
			job.JobID = jobID
			errMsg := job.Error
			if errMsg == "" {
				errMsg = "scrape job " + job.Status
			}
			return job, &JobFailedError{Message: errMsg, JobID: jobID, Status: job.Status}
		}

		if time.Now().Add(o.PollInterval).After(deadline) {
			return nil, &TimeoutError{
				Message:        fmt.Sprintf("spidra: scrape job timed out after %s, jobId: %s", o.Timeout, jobID),
				JobID:          jobID,
				TimeoutSeconds: int64(o.Timeout.Seconds()),
			}
		}

		time.Sleep(o.PollInterval)
	}
}

// RunURL scrapes a single URL with default options and waits for completion.
func (r *ScrapeResource) RunURL(ctx context.Context, targetURL string, opts ...PollOptions) (*ScrapeJob, error) {
	return r.Run(ctx, ScrapeParams{
		URLs: []ScrapeURL{{URL: targetURL}},
	}, opts...)
}
