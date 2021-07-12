package main

import (
	"database/sql"
	"fmt"
	"mysql-exporter/collectors"
	"mysql-exporter/config"
	"mysql-exporter/handler"
	"net"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/natefinch/lumberjack"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

//TODO
// init初始化
// viper配置文件
// log 日志
// 引用mysql

func initLogger(options config.Logger) func() {
	logger := lumberjack.Logger{
		Filename:   options.FileName,
		MaxSize:    options.MaxSize,
		MaxAge:     options.MaxAge, //days
		MaxBackups: options.MaxBackups,
		Compress:   options.Compress,
	}

	logrus.SetOutput(&logger)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)

	return func() {
		logger.Close()
	}
}

func initDb(options config.MySQL) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		options.Username,
		options.Password,
		options.Host,
		options.Port,
		options.Db,
	)
	return sql.Open("mysql", dsn)
}

func initMetrics(options *config.Options, db *sql.DB) {
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "mysql_up",
			Help: "MySQL UP Info",
			ConstLabels: prometheus.Labels{
				"addr": net.JoinHostPort(
					options.MySQL.Host,
					strconv.Itoa(options.MySQL.Port),
				),
			},
		},
		func() float64 {
			if err := db.Ping(); err == nil {
				return 1
			} else {
				logrus.WithFields(logrus.Fields{
					"metric": "mysql_up",
				}).Error(err)
			}
			return 0
		},
	))

	prometheus.MustRegister(collectors.NewSlowQueriesCollector(db))
	prometheus.MustRegister(collectors.NewQpsCollector(db))
	prometheus.MustRegister(collectors.NewCommandCollector(db))
	prometheus.MustRegister(collectors.NewConnectionCollector(db))
	prometheus.MustRegister(collectors.NewTrafficCollector(db))
}
func main() {
	// 指标项
	// 指标类型/触发时间 采集api的时候触发
	// 慢查询 Counter
	// 执行查询数量 Counter
	// 等等
	options, err := config.ParseConfig("./etc/exporter.yaml")
	if err != nil {
		logrus.Error(err)
	}
	close := initLogger(options.Logger)
	defer close()

	db, err := initDb(options.MySQL)
	if err != nil && db.Ping() != nil {
		logrus.Fatal(err)
	}

	initMetrics(options, db)

	// txt, _ := bcrypt.GenerateFromPassword([]byte("123abc"), 5)
	// fmt.Println(string(txt)) // $2a$05$pJJshN/.PRj2a59KkfWFXeP3a.1L3Iq9VAj4Ny/1hxVexWAVg2Tvq

	// fmt.Println(db)
	// 暴露指标
	http.Handle("/metrics", handler.Auth(
		promhttp.Handler(),
		handler.AuthSecrets{
			options.Web.Auth.Username: options.Web.Auth.Password},
	))
	http.ListenAndServe(options.Web.Addr, nil)

}
