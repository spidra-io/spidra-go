package spidra

import "context"

// MapResource handles domain mapping and sitemap URL discovery.
type MapResource struct {
	http *httpClient
}

// Run maps a domain and returns all discovered URLs.
func (r *MapResource) Run(ctx context.Context, params MapParams) (*MapResult, error) {
	var result MapResult
	if err := r.http.post(ctx, "/map", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RunURL is a convenience helper to map a domain starting from a URL string.
func (r *MapResource) RunURL(ctx context.Context, targetURL string) (*MapResult, error) {
	return r.Run(ctx, MapParams{URL: targetURL})
}

// RunURLWithSearch maps a domain filtering links by a keyword.
func (r *MapResource) RunURLWithSearch(ctx context.Context, targetURL, search string) (*MapResult, error) {
	return r.Run(ctx, MapParams{URL: targetURL, Search: search})
}

