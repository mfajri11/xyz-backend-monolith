package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Option func(db *sql.DB)

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(db *sql.DB) {
		db.SetMaxIdleConns(maxIdleConns)
	}
}

func WithMaxOpenConns(maxOpenConns int) Option {
	return func(db *sql.DB) {
		db.SetMaxOpenConns(maxOpenConns)
	}
}

func WithMaxLifetimeConn(connMaxLifetime time.Duration) Option {
	return func(db *sql.DB) {
		db.SetConnMaxLifetime(connMaxLifetime)
	}
}

func DataSource(username, password, host string, port int, databaseName string) string {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		username, password, host, port, databaseName)
	return dataSource
}

func MustNew(source string, opts ...Option) *sql.DB {

	fmt.Println(source)
	db, err := sql.Open("mysql", source)
	if err != nil {
		panic(err)
	}

	for _, opt := range opts {
		opt(db)
	}

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			continue
		}
	}

	if err != nil {
		panic(err)
	}

	return db
}
