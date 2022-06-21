package db

import (
	"mfbo_api/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB
var DbDashboard *sqlx.DB

func DBInit() *sqlx.DB {
	db := sqlx.MustConnect("mysql", config.DBUser+":"+config.DBPassword+"@tcp("+config.DBHost+")/"+config.DBName+"?charset=latin1")
	return db
}

func DBDashboardInit() *sqlx.DB {
	db := sqlx.MustConnect("mysql", config.DBUser+":"+config.DBPassword+"@tcp("+config.DBHost+")/"+config.DBDashboardName+"?charset=latin1")
	return db
}
