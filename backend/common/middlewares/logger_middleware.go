package middlewares

import (
	"bytes"
	"common/dto"
	"common/pkg/jwt"
	"common/pkg/logger"
	"common/pkg/utils"
	"io"
	"time"

	"common/constants"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(logger logger.Logger, jwt jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var requestBodyBytes []byte

		if c.Request.Body != nil {
			requestBodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
		}

		logEntry := dto.LogEntry{
			RequestID:   c.GetString("request_id"),
			TraceID:     c.GetString("trace_id"),
			SpanID:      c.GetString("span_id"),
			Host:        c.Request.Host,
			Path:        c.Request.URL.RequestURI(),
			ClientIP:    c.ClientIP(),
			LatLong:     c.GetHeader("X-Lat-Long"),
			UserID:      c.GetString("user_id"),
			RequestSize: utils.FormatMemorySize(requestBodyBytes),
			Method:      c.Request.Method,
			RequestBody: string(requestBodyBytes),
			UserAgent:   c.Request.UserAgent(),
			Action:      constants.ActionMiddlewareStart,
			EndTime:     time.Now(),
		}
		logEntry.Action = constants.ActionMiddlewareStart
		logger.Info("Request started", logEntry)

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		endTime := time.Now()
		logEntry.EndTime = endTime
		logEntry.ResponseBody = blw.body.String()
		logEntry.StatusCode = blw.Status()

		latency := endTime.Sub(startTime)

		logEntry.Duration = latency
		logEntry.Action = constants.ActionMiddlewareEnd
		logEntry.ResponseSize = utils.FormatMemorySize(blw.body.Bytes())

		switch {
		case blw.Status() >= 500:
			logEntry.Action = constants.ActionMiddlewareError
			logger.Error("Server Error", logEntry)
		case blw.Status() >= 400:
			logEntry.Action = constants.ActionMiddlewareError
			logger.Warn("Client Error", logEntry)
		default:
			logEntry.Action = constants.ActionMiddlewareEnd
			logger.Info("Request Completed", logEntry)
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
