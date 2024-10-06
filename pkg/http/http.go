package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpClientService interface {
	Get(url string, headers map[string]string) ([]byte, error)
	Post(url string, body []byte, headers map[string]string) ([]byte, error)
	Put(url string, body []byte, headers map[string]string) ([]byte, error)
	Delete(url string, headers map[string]string) ([]byte, error)
}

type httpClientService struct {
	client        *http.Client
	retryCount    int
	globalHeaders map[string]string
	baseURL       string
}

func New(config ...Config) *httpClientService {
	cfg := defaultConfig(config...)

	return &httpClientService{
		client:        &http.Client{Timeout: cfg.Timeout},
		retryCount:    cfg.RetryCount,
		globalHeaders: cfg.GlobalHeaders,
		baseURL:       cfg.BaseURL,
	}
}

func (h *httpClientService) addHeaders(req *http.Request, headers ...map[string]string) {
	for key, value := range h.globalHeaders {
		req.Header.Add(key, value)
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Add(key, value)
		}
	}

}

func (s *httpClientService) resolveURL(endpoint string) string {
	if s.baseURL == "" {
		return endpoint
	}

	resolvedURL, err := url.JoinPath(s.baseURL, endpoint)
	if err != nil {
		fmt.Println("error", err)
		return endpoint
	}
	return resolvedURL
}

func (h *httpClientService) makeRequest(req *http.Request) ([]byte, error) {
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

func (h *httpClientService) Get(endpoint string, headers ...map[string]string) ([]byte, error) {
	url := h.resolveURL(endpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	h.addHeaders(req, headers...)

	return h.makeRequest(req)
}

func (h *httpClientService) Post(endpoint string, body []byte, headers ...map[string]string) ([]byte, error) {
	url := h.resolveURL(endpoint)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	h.addHeaders(req, headers...)
	return h.makeRequest(req)
}

func (h *httpClientService) Put(endpoint string, body []byte, headers ...map[string]string) ([]byte, error) {
	url := h.resolveURL(endpoint)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	h.addHeaders(req, headers...)
	return h.makeRequest(req)
}

func (h *httpClientService) Delete(endpoint string, headers ...map[string]string) ([]byte, error) {
	url := h.resolveURL(endpoint)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	h.addHeaders(req, headers...)
	return h.makeRequest(req)
}
