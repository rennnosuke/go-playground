package conn

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type mySQLConfig struct {
	userName string
	password string
	dbName   string
}

func (c mySQLConfig) getDataSourceName() string {
	return c.userName + ":" + c.password + "@/" + c.dbName
}

func open(_ context.Context, conf mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	return db, nil
}
