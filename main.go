package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status_code"},
	)

	activeRequestGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active HTTP requests",
		},
	)

	requestDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
		},
		[]string{"method", "route", "status_code"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(activeRequestGauge)
	prometheus.MustRegister(requestDurationHistogram)
}

func metricsMiddleware(next http.HandlerFunc, routeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activeRequestGauge.Inc()
		defer activeRequestGauge.Dec()

		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next(rw, r)

		duration := time.Since(start).Seconds()
		statusCode := fmt.Sprintf("%d", rw.statusCode)

		requestCounter.WithLabelValues(r.Method, routeName, statusCode).Inc()
		requestDurationHistogram.WithLabelValues(r.Method, routeName, statusCode).Observe(duration)
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func mathHandler(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 1000000000; i++ {
		_ = rand.Float64() * 10
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "done!",
	})
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user fetched!",
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func main() {
	http.HandleFunc("/math", metricsMiddleware(mathHandler, "/math"))
	http.HandleFunc("/user", metricsMiddleware(userHandler, "/user"))
	http.HandleFunc("/health", metricsMiddleware(healthHandler, "/health"))

	http.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

