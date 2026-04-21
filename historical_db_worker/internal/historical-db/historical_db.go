package historicaldb

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	redisdb "histo-db/internal/redis_db"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func generateConnString() (c string) {
	dbUser := os.Getenv("POSTGRES_DB_USER")
	dbPassword := os.Getenv("POSTGRES_DB_PASSWORD")
	pgHost := os.Getenv("POSTGRES_DB_HOST")
	dbName := os.Getenv("POSTGRES_DB_NAME")
	c = ("postgres://" + dbUser + ":" + dbPassword + "@" + pgHost + ":5432/" + dbName + "?sslmode=disable")
	return
}

func ConnectToHistoricalDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file, recovering variables from environment")
	}

	connStr := generateConnString()

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
