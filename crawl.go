package spidra

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// CrawlResource handles deep site crawling and multi-page extraction jobs.
type CrawlResource struct {
	http *httpClient
}

// Stats returns the total number of crawl jobs for your account.
func (r *CrawlResource) Stats(ctx context.Context) (*CrawlStats, error) {
	var stats CrawlStats
	if err := r.http.get(ctx, "/crawl/stats", nil, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// History lists past crawl jobs, newest first.
func (r *CrawlResource) History(ctx context.Context, page, limit int) (*CrawlHistory, error) {
	var history CrawlHistory
	query := map[string]string{}
	if page > 0 {
		query["page"] = strconv.Itoa(page)
	}
	if limit > 0 {
		query["limit"] = strconv.Itoa(limit)
	}
	if err := r.http.get(ctx, "/crawl/history", query, &history); err != nil {
		return nil, err
	}
	return &history, nil
}

// Submit queues a crawl job and returns immediately with a JobID.
func (r *CrawlResource) Submit(ctx context.Context, params CrawlParams) (*CrawlJob, error) {
	var job CrawlJob
	if err := r.http.post(ctx, "/crawl", params, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// SubmitURL is a convenience helper to crawl a seed URL with default instructions.
func (r *CrawlResource) SubmitURL(ctx context.Context, baseURL string) (*CrawlJob, error) {
	return r.Submit(ctx, CrawlParams{BaseURL: baseURL})
}

// Get returns the current status and progress of a crawl job.
func (r *CrawlResource) Get(ctx context.Context, jobID string) (*CrawlJob, error) {
	var job CrawlJob
	if err := r.http.get(ctx, "/crawl/"+jobID, nil, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// Pages returns all pages collected by a completed crawl, with signed download URLs.
func (r *CrawlResource) Pages(ctx context.Context, jobID string) (*CrawlPages, error) {
	var pages CrawlPages
	if err := r.http.get(ctx, "/crawl/"+jobID+"/pages", nil, &pages); err != nil {
		return nil, err
	}
	return &pages, nil
}

// Extract re-runs AI extraction on a completed crawl with a new instruction.
// Does not re-crawl — only reprocesses already collected pages.
func (r *CrawlResource) Extract(ctx context.Context, jobID, transformInstruction string) (*CrawlExtract, error) {
	var result CrawlExtract
	err := r.http.post(ctx, "/crawl/"+jobID+"/extract", map[string]any{
		"transformInstruction": transformInstruction,
	}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Cancel stops an active crawl job and refunds credits for uncrawled pages.
func (r *CrawlResource) Cancel(ctx context.Context, jobID string) (*CrawlCancelResult, error) {
	var result CrawlCancelResult
	if err := r.http.delete(ctx, "/crawl/"+jobID, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Report retrieves a diagnostic report of skipped URLs and crawl errors.
func (r *CrawlResource) Report(ctx context.Context, jobID string) (*CrawlReportResponse, error) {
	var resp CrawlReportResponse
	if err := r.http.get(ctx, "/crawl/"+jobID+"/report", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RetryPage retries crawling a specific failed page in a crawl job.
func (r *CrawlResource) RetryPage(ctx context.Context, jobID, pageID string) (*CrawlRetryPageResponse, error) {
	var resp CrawlRetryPageResponse
	if err := r.http.post(ctx, "/crawl/"+jobID+"/pages/"+pageID+"/retry", map[string]any{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Run submits a crawl job and polls until it completes.
func (r *CrawlResource) Run(ctx context.Context, params CrawlParams, opts ...PollOptions) (*CrawlJob, error) {
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
			errMsg = "crawl job " + queued.Status
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
				errMsg = "crawl job " + job.Status
			}
			return job, &JobFailedError{Message: errMsg, JobID: jobID, Status: job.Status}
		}

		if time.Now().Add(o.PollInterval).After(deadline) {
			return nil, &TimeoutError{
				Message:        fmt.Sprintf("spidra: crawl job timed out after %s, jobId: %s", o.Timeout, jobID),
				JobID:          jobID,
				TimeoutSeconds: int64(o.Timeout.Seconds()),
			}
		}

		time.Sleep(o.PollInterval)
	}
}

// RunURL crawls a seed URL and polls until completion.
func (r *CrawlResource) RunURL(ctx context.Context, baseURL string, opts ...PollOptions) (*CrawlJob, error) {
	return r.Run(ctx, CrawlParams{BaseURL: baseURL}, opts...)
}
