package store

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4/pgxpool"
	"miikka.xyz/rs/config"
)

// DB is wrapper for database instance
type DB struct {
	*pgxpool.Pool
}

// NewDB returns database instance
func NewDB(config *config.Config) *DB {
	connStr := "postgresql://" + config.DBUser + ":" + config.DBPassword + "@" + config.DBURL + "/" + config.DBName
	log.Info("connection string: ", connStr)
	conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("database connection established")
	return &DB{conn}
}
