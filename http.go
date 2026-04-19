package spidra

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type httpClient struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func (h *httpClient) get(ctx context.Context, path string, query map[string]string, out any) error {
	return h.do(ctx, http.MethodGet, path, query, nil, out)
}

func (h *httpClient) post(ctx context.Context, path string, body, out any) error {
	return h.do(ctx, http.MethodPost, path, nil, body, out)
}

func (h *httpClient) delete(ctx context.Context, path string, out any) error {
	return h.do(ctx, http.MethodDelete, path, nil, nil, out)
}

func (h *httpClient) do(ctx context.Context, method, path string, query map[string]string, body, out any) error {
	u := strings.TrimRight(h.baseURL, "/") + "/" + strings.TrimLeft(path, "/")

	if len(query) > 0 {
		params := url.Values{}
		for k, v := range query {
			params.Set(k, v)
		}
		u += "?" + params.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", h.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := h.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return mapError(resp.StatusCode, data)
	}

	if out != nil {
		return json.Unmarshal(data, out)
	}
	return nil
}
