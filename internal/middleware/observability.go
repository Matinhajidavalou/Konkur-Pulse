package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"konkur-pulse/internal/metrics"

	"github.com/google/uuid"
)

// RequestIDKey represents the context key for injecting request IDs.
type RequestIDKey string

const ReqIDKey RequestIDKey = "request_id"

// Observability acts as a middleware wrapping the HTTP router.
// It handles structured logging, tracing (Request ID injection), and metrics collection.
func Observability(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Inject unique request ID for distributed tracing
			reqID := uuid.New().String()
			ctx := context.WithValue(r.Context(), ReqIDKey, reqID)
			r = r.WithContext(ctx)

			logger.Info("Incoming request",
				slog.String("request_id", reqID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			next.ServeHTTP(w, r)

			// Record latency and request counts
			duration := time.Since(start)
			metrics.HttpRequestsTotal.WithLabelValues(r.URL.Path, r.Method).Inc()
			metrics.HttpRequestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration.Seconds())

			logger.Info("Request completed",
				slog.String("request_id", reqID),
				slog.Duration("duration", duration),
			)
		})
	}
}
