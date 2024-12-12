package http_client

import (
	"net/http"
	"time"
)

type Config struct {
	Timeout       time.Duration
	RetryCount    int
	BaseURL       string
	GlobalHeaders map[string]string
	Interceptor   http.RoundTripper
}

var ConfigDefault = Config{
	Timeout: 30 * time.Second,
}

// defaultConfig returns a default config for http config.
func defaultConfig(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	return cfg
}
