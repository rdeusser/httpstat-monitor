package httpstat

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	urlUp = promauto.NewGauge(prometheus.GaugeOpts{})
)
