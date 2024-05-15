package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
}

func NewMetrics(
	reg prometheus.Registerer,
) *Metrics {
	m := &Metrics{}
	return m
}
