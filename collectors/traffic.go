package collectors

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type TrafficCollector struct {
	mysqlCollector
	desc *prometheus.Desc
}

func NewTrafficCollector(db *sql.DB) *TrafficCollector {
	return &TrafficCollector{
		mysqlCollector: mysqlCollector{db},
		desc: prometheus.NewDesc(
			"mysql_global_status_traffic",
			"MySQL global status traffic",
			[]string{"direction"},
			nil,
		),
	}
}

func (c *TrafficCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- c.desc
}

func (c *TrafficCollector) Collect(metrics chan<- prometheus.Metric) {
	received := c.status("Bytes_received")
	sent := c.status("Bytes_sent")
	metrics <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, received, "received")
	metrics <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, sent, "sent")
}
