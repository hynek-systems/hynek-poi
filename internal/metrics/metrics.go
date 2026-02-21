package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hynek_poi_requests_total",
			Help: "Total number of search requests",
		},
		[]string{"endpoint"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hynek_poi_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	CacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hynek_poi_cache_hits_total",
			Help: "Total cache hits",
		},
	)

	CacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "hynek_poi_cache_misses_total",
			Help: "Total cache misses",
		},
	)

	ProviderDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hynek_poi_provider_duration_seconds",
			Help:    "Provider request duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"provider"},
	)

	ProviderErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hynek_poi_provider_errors_total",
			Help: "Provider errors",
		},
		[]string{"provider"},
	)
)

func Register() {

	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(CacheHits)
	prometheus.MustRegister(CacheMisses)
	prometheus.MustRegister(ProviderDuration)
	prometheus.MustRegister(ProviderErrors)
}
