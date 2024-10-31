package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

func NewConnection(host, user, password, dbname string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, dbname)
	return sqlx.Connect("mysql", dsn)
}