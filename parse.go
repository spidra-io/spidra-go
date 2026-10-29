package spidra

import (
	"context"
	"fmt"
	"time"
)

// ParseResource handles document (PDF, Office) and media file parsing.
type ParseResource struct {
	http *httpClient
}

// Submit queues a document parsing job and returns immediately with a ParseJob.
func (r *ParseResource) Submit(ctx context.Context, params ParseParams) (*ParseJob, error) {
	var job ParseJob
	if err := r.http.post(ctx, "/parse", params, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// SubmitURL is a convenience helper to submit a document URL for parsing.
func (r *ParseResource) SubmitURL(ctx context.Context, docURL string) (*ParseJob, error) {
	return r.Submit(ctx, ParseParams{URL: docURL})
}

// Get retrieves the status and extracted content of a parse job.
func (r *ParseResource) Get(ctx context.Context, jobID string) (*ParseJob, error) {
	var job ParseJob
	if err := r.http.get(ctx, "/parse/"+jobID, nil, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// Run submits a document parse request and polls until it completes.
func (r *ParseResource) Run(ctx context.Context, params ParseParams, opts ...PollOptions) (*ParseJob, error) {
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
			errMsg = "parse job " + job.Status
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
				errMsg = "parse job " + current.Status
			}
			return current, &JobFailedError{Message: errMsg, JobID: jobID, Status: current.Status}
		}

		if time.Now().Add(o.PollInterval).After(deadline) {
			return nil, &TimeoutError{
				Message:        fmt.Sprintf("spidra: parse job timed out after %s, jobId: %s", o.Timeout, jobID),
				JobID:          jobID,
				TimeoutSeconds: int64(o.Timeout.Seconds()),
			}
		}

		time.Sleep(o.PollInterval)
	}
}

// RunURL submits a document URL and waits for parsed results.
func (r *ParseResource) RunURL(ctx context.Context, docURL string, opts ...PollOptions) (*ParseJob, error) {
	return r.Run(ctx, ParseParams{URL: docURL}, opts...)
}

// Estimate calculates credit cost estimates for parsing a document.
func (r *ParseResource) Estimate(ctx context.Context, params ParseParams) (*ParseEstimate, error) {
	var estimate ParseEstimate
	if err := r.http.post(ctx, "/parse/estimate", params, &estimate); err != nil {
		return nil, err
	}
	return &estimate, nil
}

// EstimateURL is a convenience helper to estimate credits for parsing a URL.
func (r *ParseResource) EstimateURL(ctx context.Context, docURL string) (*ParseEstimate, error) {
	return r.Estimate(ctx, ParseParams{URL: docURL})
}

// DeleteFile deletes an uploaded parse file by its ID.
func (r *ParseResource) DeleteFile(ctx context.Context, fileID string) error {
	return r.http.delete(ctx, "/parse/files/"+fileID, nil)
}

// SubmitBatch queues a batch of files or URLs for parsing.
func (r *ParseResource) SubmitBatch(ctx context.Context, batchParams map[string]any) (*ParseBatch, error) {
	var batch ParseBatch
	if err := r.http.post(ctx, "/parse/batch", batchParams, &batch); err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetBatch retrieves the current status of a parse batch.
func (r *ParseResource) GetBatch(ctx context.Context, batchID string) (*ParseBatch, error) {
	var batch ParseBatch
	if err := r.http.get(ctx, "/parse/batch/"+batchID, nil, &batch); err != nil {
		return nil, err
	}
	return &batch, nil
}

