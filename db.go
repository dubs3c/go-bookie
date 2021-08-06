package gobookie

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DBInit - Connects to database
func DBInit() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("host=127.0.0.1 port=5432 user=bookie dbname=bookie password=bookie sslmode=disable")

	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}

	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return dbpool, err
}
