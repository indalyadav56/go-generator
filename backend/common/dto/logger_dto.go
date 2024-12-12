package dto

import (
	"common/constants"
	"time"
)

type LogEntry struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	RequestID    string `json:"request_id"`
	TraceID      string `json:"trace_id"`
	SpanID       string `json:"span_id"`
	ParentSpanID string `json:"parent_span_id"`

	UserID string `json:"user_id"`

	Method string `json:"method"`

	Host string `json:"host"`
	Path string `json:"path"`

	ClientIP   string        `json:"client_ip"`
	LatLong    string        `json:"lat_long"`
	StatusCode int           `json:"status_code"`
	Duration   time.Duration `json:"duration"`

	UserAgent string `json:"user_agent"`

	RequestSize   string      `json:"request_size"`
	RequestHeader interface{} `json:"request_header"`
	RequestBody   interface{} `json:"request_body"`

	ResponseSize string `json:"response_size"`
	ResponseBody string `json:"response_body"`

	Action constants.LogAction `json:"action"`
}
