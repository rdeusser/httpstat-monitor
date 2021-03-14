package httpstat

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	urlUp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sample_external_url_up",
		Help: "Is the URL up?",
	})

	urlResponseMS = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "sample_external_url_response_ms",
		Help: "URL response time in milliseconds",
	})
)
