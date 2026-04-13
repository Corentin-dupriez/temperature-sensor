package historicaldb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type tempReading struct {
	temperature float64
	huimidity   float64
}

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
	fmt.Println("Connected to Postgres")
	return db
}

func addTempReading(db *sql.DB, t tempReading) (int64, error) {
	result, err := db.Exec("INSERT INTO readings (temperature, humidity) VALUES(?,?)", t.temperature, t.huimidity)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id, nil
}
