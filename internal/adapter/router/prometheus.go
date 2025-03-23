package router

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)
	responseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "Response time in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	// Register custom metrics
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(responseTime)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(responseTime.WithLabelValues(r.Method, r.URL.Path))
		defer timer.ObserveDuration()

		// Increment the request counter
		requestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()

		next.ServeHTTP(w, r)
	})
}
