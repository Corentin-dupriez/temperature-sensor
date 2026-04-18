package historicaldb

import (
	"database/sql"
	"fmt"
	"log"

	redisdb "histo-db/internal/redis_db"

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
	fmt.Println("Connected to Postgres")
	res, err := db.Exec("CREATE TABLE IF NOT EXISTS readings (id serial, temperature float, humidity float, time_reading timestamp)")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	return db
}

func AddTempReading(db *sql.DB, t redisdb.TempReading) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO readings (temperature, humidity, time_reading) VALUES ($1, $2, $3) RETURNING id", t.Temperature, t.Humidity, t.TimeReading).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("AddTempReading: %v", err)
	}
	return id, nil
}
