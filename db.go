package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
	connStr := "user=postgres password= dbname=mogu10botdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return db, err
}
