package gobookie

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

// DBInit - Connects to database
func DBInit() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig("host=postgres port=5432 user=bookie dbname=bookie password=bookie sslmode=disable")

	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}

	looger := &log.Logger{
		Out:          os.Stderr,
		Formatter:    new(log.JSONFormatter),
		Hooks:        make(log.LevelHooks),
		Level:        log.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	config.ConnConfig.Logger = logrusadapter.NewLogger(looger)

	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return dbpool, err
}
