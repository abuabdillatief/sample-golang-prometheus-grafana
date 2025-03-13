package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels = []string{"path", "status"}

	reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of requests processed by the server.",
		},
		labels,
	)

	errCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_errors_total",
			Help: "Total number of error requests processed by the server.",
		},
		labels,
	)

	durrCount = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Histogram of the response latency (seconds) of the server.",
			Buckets: prometheus.DefBuckets,
		},
		labels,
	)
)

func PrometheusInit() {
	prometheus.MustRegister(reqCount)
	prometheus.MustRegister(errCount)
	prometheus.MustRegister(durrCount)
}

func TrackMetrics(c *gin.Context) {
	path := c.Request.URL.Path
	start := time.Now()

	c.Next()

	status := c.Writer.Status()
	duration := time.Since(start).Seconds()

	statusText := http.StatusText(status)
	reqCount.WithLabelValues(path, statusText).Inc()
	durrCount.WithLabelValues(path, statusText).Observe(duration)

	if status >= 400 {
		errCount.WithLabelValues(path, statusText).Inc()
	}
}
