package spidra

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// CrawlResource handles crawl jobs.
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
	err := r.http.get(ctx, "/crawl/history", map[string]string{
		"page":  strconv.Itoa(page),
		"limit": strconv.Itoa(limit),
	}, &history)
	if err != nil {
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
			return nil, fmt.Errorf("spidra: crawl job timed out after %s, jobId: %s", o.Timeout, jobID)
		}

		time.Sleep(o.PollInterval)
	}
}
