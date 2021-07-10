package collectors

import (
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type SlowQueriesCollector struct {
	db   *sql.DB
	desc *prometheus.Desc
}

func NewSlowQueriesCollector(db *sql.DB) *SlowQueriesCollector {
	return &SlowQueriesCollector{
		db: db,
		desc: prometheus.NewDesc(
			"mysql_global_status_slow_queries",
			"Mysql global status slow queries",
			nil,
			nil,
		),
	}
}

func (c *SlowQueriesCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- c.desc
}

func (c *SlowQueriesCollector) Collect(metrics chan<- prometheus.Metric) {
	var (
		name  string
		count float64
	)
	c.db.QueryRow("show global status where variable_name=?", "slow_queries").Scan(&name, &count)
	metrics <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, count)
}

type QpsCollector struct {
	mysqlCollector
	desc *prometheus.Desc
}

func NewQpsCollector(db *sql.DB) *QpsCollector {
	return &QpsCollector{
		mysqlCollector: mysqlCollector{db},
		desc: prometheus.NewDesc(
			"mysql_global_status_qps",
			"Mysql global status qps",
			nil,
			nil,
		),
	}
}

func (c *QpsCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- c.desc
}

func (c *QpsCollector) Collect(metrics chan<- prometheus.Metric) {
	var (
		name  string
		count float64
	)
	err := c.db.QueryRow("show global status where variable_name=?", "queries").Scan(&name, &count)
	if err != nil {
		fmt.Println(err)
	}
	metrics <- prometheus.MustNewConstMetric(c.desc, prometheus.CounterValue, count)
}
