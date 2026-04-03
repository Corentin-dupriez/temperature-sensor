package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectToHistoricalDB() *sql.DB {
	connStr := "postgres://postgres:example@localhost:5432/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
