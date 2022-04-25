package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDB(сonnStr string) (dbr *sql.DB, err error) {
	if db != nil {
		return db, nil
	}

	dbr, err = sql.Open("postgres", сonnStr)
	return dbr, err
}
