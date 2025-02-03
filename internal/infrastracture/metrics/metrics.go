package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	GRPCRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	GRPCRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC request",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	ExternalAPIRequests = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "external_api_requests_total",
			Help:    "Total requests to Gurantex API",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)

	ExternalAPIDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "external_api_duration_seconds",
			Help:    "Duration of Garantex API requests",
			Buckets: prometheus.DefBuckets,
		},
	)

	DBRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_requests_total",
			Help: "Total database operations",
		},
		[]string{"operation"},
	)

	DBDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_duration_seconds",
			Help:    "Database operations duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	HealthcheckStatus = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "healthcheck_status",
			Help: "Service healthcheck status (1 = healthy, 2 = unhealthy)",
		},
	)
)

func InitMetrics(port string) {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":"+port, nil)
	}()
}
