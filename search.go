package spidra

import (
	"context"
	"fmt"
	"time"
)

// SearchResource handles web, news, and multimodal search jobs.
type SearchResource struct {
	http *httpClient
}

// Submit queues a search job and returns immediately with a SearchJob.
func (r *SearchResource) Submit(ctx context.Context, params SearchParams) (*SearchJob, error) {
	var job SearchJob
	if err := r.http.post(ctx, "/search", params, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// SubmitQuery is a convenience helper to submit a simple query string.
func (r *SearchResource) SubmitQuery(ctx context.Context, query string) (*SearchJob, error) {
	return r.Submit(ctx, SearchParams{Query: query})
}

// Get retrieves the current status and results of a search job.
func (r *SearchResource) Get(ctx context.Context, jobID string) (*SearchJob, error) {
	var job SearchJob
	if err := r.http.get(ctx, "/search/"+jobID, nil, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// Cancel stops an in-progress search job.
func (r *SearchResource) Cancel(ctx context.Context, jobID string) error {
	return r.http.delete(ctx, "/search/"+jobID, nil)
}

// Run executes a search query and waits for results.
func (r *SearchResource) Run(ctx context.Context, params SearchParams, opts ...PollOptions) (*SearchJob, error) {
	o := PollOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	o = o.withDefaults()

	params.Wait = true
	job, err := r.Submit(ctx, params)
	if err != nil {
		return nil, err
	}

	jobID := job.JobID

	switch job.Status {
	case "completed":
		return job, nil
	case "failed", "error", "cancelled":
		errMsg := job.Error
		if errMsg == "" {
			errMsg = "search job " + job.Status
		}
		return job, &JobFailedError{Message: errMsg, JobID: jobID, Status: job.Status}
	}

	if jobID == "" {
		return job, nil
	}

	deadline := time.Now().Add(o.Timeout)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		current, err := r.Get(ctx, jobID)
		if err != nil {
			return nil, err
		}

		switch current.Status {
		case "completed":
			current.JobID = jobID
			return current, nil
		case "failed", "error", "cancelled":
			current.JobID = jobID
			errMsg := current.Error
			if errMsg == "" {
				errMsg = "search job " + current.Status
			}
			return current, &JobFailedError{Message: errMsg, JobID: jobID, Status: current.Status}
		}

		if time.Now().Add(o.PollInterval).After(deadline) {
			return nil, &TimeoutError{
				Message:        fmt.Sprintf("spidra: search job timed out after %s, jobId: %s", o.Timeout, jobID),
				JobID:          jobID,
				TimeoutSeconds: int64(o.Timeout.Seconds()),
			}
		}

		time.Sleep(o.PollInterval)
	}
}

// RunQuery searches with a query string and waits for results.
func (r *SearchResource) RunQuery(ctx context.Context, query string, opts ...PollOptions) (*SearchJob, error) {
	return r.Run(ctx, SearchParams{Query: query}, opts...)
}

