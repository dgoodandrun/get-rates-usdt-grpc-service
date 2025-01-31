package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of request received",
		},
		[]string{"method", "endpoint"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of http requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	CacheDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_request_duration_seconds",
			Help:    "Duration of cache request",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	DBDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_duration_request_seconds",
			Help:    "Duration of DB requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	APIDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_api_request_duration_seconds",
			Help:    "Duration of external API requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	AuthDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_request_duration_seconds",
			Help:    "Duration of auth requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func Init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(CacheDuration)
	prometheus.MustRegister(DBDuration)
	prometheus.MustRegister(APIDuration)
	prometheus.MustRegister(AuthDuration)
}
