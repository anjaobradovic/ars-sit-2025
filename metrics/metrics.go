package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Ukupan broj HTTP zahteva
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "config_service_http_requests_total",
			Help: "Ukupan broj HTTP zahteva za configuration service",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Uspešni zahtevi (2xx, 3xx)
	HttpRequestsSuccessful = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "config_service_http_requests_successful_total",
			Help: "Broj uspešnih HTTP zahteva (2xx, 3xx)",
		},
		[]string{"method", "endpoint"},
	)

	// Neuspešni zahtevi (4xx, 5xx)
	HttpRequestsFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "config_service_http_requests_failed_total",
			Help: "Broj neuspešnih HTTP zahteva (4xx, 5xx)",
		},
		[]string{"method", "endpoint", "status_class"},
	)

	// Histogram vremena odgovora
	HttpResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "config_service_http_response_duration_seconds",
			Help:    "Vreme odgovora HTTP zahteva u sekundama",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Trenutno aktivni zahtevi
	HttpRequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "config_service_http_requests_in_flight",
			Help: "Broj trenutno aktivnih HTTP zahteva",
		},
		[]string{"method", "endpoint"},
	)

	registry = prometheus.NewRegistry()
)

func init() {
	registry.MustRegister(
		HttpRequestsTotal,
		HttpRequestsSuccessful,
		HttpRequestsFailed,
		HttpResponseDuration,
		HttpRequestsInFlight,
	)
}

func MetricsHandler() http.Handler {
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
