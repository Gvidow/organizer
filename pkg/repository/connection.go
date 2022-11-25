package repository

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const MaxOpenConns = 100_000

func ConnectDB() (*sql.DB, error) {
	dsn := "root:admin@tcp(localhost:3306)/organizer?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(MaxOpenConns)
	err = db.Ping()
	return db, err
}
