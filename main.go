package main

import (
	"database/sql"
	"log"
	"mysql-exporter/collectors"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 指标项
	// 指标类型/触发时间 采集api的时候触发
	// 慢查询 Counter
	// 执行查询数量 Counter
	// 等等

	dsn := "root:kubernetes@2020@tcp(121.36.50.10:3306)/mysql"
	mysqlAddr := "121.36.50.10:3306"
	addr := ":9002"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(db)

	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "mysql_up",
			Help:        "Mysql UP Info",
			ConstLabels: prometheus.Labels{"addr": mysqlAddr},
		},
		func() float64 {
			if err := db.Ping(); err == nil {
				return 1
			}
			return 0
		},
	))

	prometheus.MustRegister(collectors.NewSlowQueriesCollector(db))
	prometheus.MustRegister(collectors.NewQpsCollector(db))

	// 暴露指标
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(addr, nil)

}
