package collectors

import (
	"database/sql"
)

type mysqlCollector struct {
	db *sql.DB
}

func (c *mysqlCollector) status(name string) float64 {
	sql := "show global status where variable_name=?"
	var (
		vname string
		rs    float64
	)
	err := c.db.QueryRow(sql, name).Scan(&vname, &rs)
	if err != nil {
		return rs
	}
	return rs
}

func (c *mysqlCollector) variables(name string) float64 {
	sql := "show global variables where variable_name=?"
	var (
		vname string
		rs    float64
	)
	err := c.db.QueryRow(sql, name).Scan(&vname, &rs)
	if err != nil {
		return rs
	}
	return rs
}
