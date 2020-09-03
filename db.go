package gobookie

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DBInit - Connects to database
func DBInit() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), "host=127.0.0.1 port=5432 user=bookie dbname=bookie password=bookie sslmode=disable")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	//defer dbpool.Close()
	return dbpool, err
}
