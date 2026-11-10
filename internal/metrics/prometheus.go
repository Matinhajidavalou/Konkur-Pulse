package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HttpRequestsTotal tracks the total number of processed HTTP requests.
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pulse_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method"},
	)

	// HttpRequestDuration tracks the latency distribution of HTTP requests.
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pulse_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	// UpstreamErrorsTotal tracks errors originating from external APIs (e.g., Maze).
	UpstreamErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pulse_upstream_errors_total",
			Help: "Total number of upstream errors and retries",
		},
		[]string{"service", "error_type"},
	)
)
