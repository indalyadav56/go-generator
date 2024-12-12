package http_client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpClient interface {
	Post(ctx context.Context, endpoint string, body []byte, headers ...map[string]string) (*Response, error)
	Get(ctx context.Context, endpoint string, headers ...map[string]string) (*Response, error)
}

type httpClient struct {
	client        *http.Client
	retryCount    int
	globalHeaders map[string]string
	baseURL       string
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type RequestError struct {
	StatusCode int
	URL        string
	Method     string
	Response   []byte
	Err        error
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("request failed: method=%s, url=%s, status=%d, error=%v",
		e.Method, e.URL, e.StatusCode, e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

func New(config ...Config) *httpClient {
	cfg := defaultConfig(config...)

	var tr http.RoundTripper
	if cfg.Interceptor == nil {
		tr = http.DefaultTransport
	} else {
		tr = cfg.Interceptor
	}

	return &httpClient{
		client: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: tr,
		},
		retryCount:    cfg.RetryCount,
		globalHeaders: cfg.GlobalHeaders,
		baseURL:       cfg.BaseURL,
	}
}

func (h *httpClient) addHeaders(req *http.Request, headers ...map[string]string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add global headers
	for key, value := range h.globalHeaders {
		req.Header.Set(key, value)
	}

	// Add request-specific headers
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}
}

func (h *httpClient) resolveURL(endpoint string) (string, error) {
	if h.baseURL == "" {
		return endpoint, nil
	}

	resolvedURL, err := url.JoinPath(h.baseURL, endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to resolve URL: %w", err)
	}
	return resolvedURL, nil
}

func (h *httpClient) makeRequest(ctx context.Context, req *http.Request) (*Response, error) {
	req = req.WithContext(ctx)

	resp, err := h.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("request canceled or timed out: %w", ctx.Err())
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		// Drain body before closing to enable connection reuse
		if resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}

	return response, nil
}

func (h *httpClient) Post(ctx context.Context, endpoint string, body []byte, headers ...map[string]string) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	url, err := h.resolveURL(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	h.addHeaders(req, headers...)
	return h.makeRequest(ctx, req)
}

func (h *httpClient) Get(ctx context.Context, endpoint string, headers ...map[string]string) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	url, err := h.resolveURL(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	h.addHeaders(req, headers...)
	return h.makeRequest(ctx, req)
}

func (h *httpClient) Put(ctx context.Context, endpoint string, body []byte, headers ...map[string]string) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	url, err := h.resolveURL(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	h.addHeaders(req, headers...)
	return h.makeRequest(ctx, req)
}

func (h *httpClient) Patch(ctx context.Context, endpoint string, body []byte, headers ...map[string]string) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	url, err := h.resolveURL(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	h.addHeaders(req, headers...)
	return h.makeRequest(ctx, req)
}

func (h *httpClient) Delete(ctx context.Context, endpoint string, headers ...map[string]string) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	url, err := h.resolveURL(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	h.addHeaders(req, headers...)
	return h.makeRequest(ctx, req)
}
