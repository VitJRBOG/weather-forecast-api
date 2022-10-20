package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Postgres driver
)

func NewConnection(dsn string) *sql.DB {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("unable connect to database: %s\n", err.Error())
	}

	return dbConn
}
