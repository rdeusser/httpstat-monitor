package httpstat

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	urlUp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_external_url_up",
		Help: "Is the URL up?",
	}, []string{"url"})

	urlResponseMS = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "sample_external_url_response_ms",
		Help:    "URL response time in milliseconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"url"})
)
