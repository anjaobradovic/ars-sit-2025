package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anjaobradovic/ars-sit-2025/metrics"
	"github.com/gorilla/mux"
)

// responseWriter wrapper
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriter) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

// Izvlači pattern rute
func getEndpointPattern(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if route == nil {
		return "unknown"
	}

	pathTemplate, err := route.GetPathTemplate()
	if err != nil {
		return "unknown"
	}

	return pathTemplate
}

// Provera uspešnosti (2xx, 3xx)
func isSuccessfulStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 400
}

// Status klasa (samo 4xx / 5xx koristimo)
func getStatusClass(statusCode int) string {
	if statusCode >= 400 && statusCode < 500 {
		return "4xx"
	}
	if statusCode >= 500 && statusCode < 600 {
		return "5xx"
	}
	return "unknown"
}

// MetricsMiddleware beleži metrike za svaki HTTP zahtev
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		method := r.Method
		endpoint := getEndpointPattern(r)

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// In-flight
		metrics.HttpRequestsInFlight.
			WithLabelValues(method, endpoint).
			Inc()
		defer metrics.HttpRequestsInFlight.
			WithLabelValues(method, endpoint).
			Dec()

		// Obrada zahteva
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		statusCode := rw.statusCode
		statusCodeStr := strconv.Itoa(statusCode)

		// Total requests
		metrics.HttpRequestsTotal.
			WithLabelValues(method, endpoint, statusCodeStr).
			Inc()

		// Response time
		metrics.HttpResponseDuration.
			WithLabelValues(method, endpoint).
			Observe(duration)

		// Success / failure
		if isSuccessfulStatusCode(statusCode) {
			metrics.HttpRequestsSuccessful.
				WithLabelValues(method, endpoint).
				Inc()
		} else {
			metrics.HttpRequestsFailed.
				WithLabelValues(method, endpoint, getStatusClass(statusCode)).
				Inc()
		}
	})
}
