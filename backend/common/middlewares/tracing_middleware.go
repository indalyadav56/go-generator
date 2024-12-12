package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var ContextKey string

func TracingMiddleware(tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 2. Extract existing trace context from headers
		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// 3. Start the root span with extracted context
		ctx, span := tracer.Start(ctx, "http.request",
			trace.WithAttributes(
				attribute.String("request_id", requestID),
				attribute.String("path", c.FullPath()),
				attribute.String("method", c.Request.Method),
				attribute.String("user_agent", c.Request.UserAgent()),
				attribute.String("remote_addr", c.ClientIP()),
			),
		)
		defer span.End()

		// 4. Get trace information from span
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		// 5. Store values in both Gin context and request context
		c.Set("request_id", requestID)
		c.Set("trace_id", traceID)
		c.Set("span_id", spanID)

		ctx = context.WithValue(ctx, "request_id", requestID)
		ctx = context.WithValue(ctx, "trace_id", traceID)

		// 6. Set response headers
		c.Header("X-Trace-ID", traceID)
		c.Header("X-Request-ID", requestID)

		// 7. Update request context
		c.Request = c.Request.WithContext(ctx)

		// 8. Process request
		c.Next()

		// 9. Add response information to span
		span.SetAttributes(
			attribute.Int("status_code", c.Writer.Status()),
			attribute.Int("response_size", c.Writer.Size()),
		)

		// 10. Record error if present
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.Bool("error", true))
			span.SetAttributes(attribute.String("error.message", c.Errors.String()))
		}
	}
}
