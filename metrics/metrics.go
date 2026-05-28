package metrics

import "github.com/prometheus/client_golang/prometheus"

var RequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gateway_request_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "path", "status"},
)

var RequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "gateway_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10, 30},
	},
	[]string{"method", "path"},
)

var CacheHitsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gateway_cache_hits_total",
		Help: "Total number of cache hits and misses",
	},
	[]string{"status"}, // "hit" or "miss"
)

var ModelRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gateway_model_requests_total",
		Help: "Total number of requests per model and provider",
	},
	[]string{"model", "provider"},
)

var RateLimitedTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gateway_rate_limited_total",
		Help: "Total number of rate limited requests",
	},
	[]string{"user_id"},
)

func Register() {
	prometheus.MustRegister(
		RequestsTotal,
		RequestDuration,
		CacheHitsTotal,
		ModelRequestsTotal,
		RateLimitedTotal,
	)
}
