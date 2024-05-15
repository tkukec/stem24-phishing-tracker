package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	QueueingDuration   *prometheus.GaugeVec
	NrOfLoggedInAgents prometheus.GaugeFunc
	NrOfQueuedItems    prometheus.GaugeFunc
}

func NewMetrics(
	reg prometheus.Registerer,
) *Metrics {
	m := &Metrics{}
	return m
}
