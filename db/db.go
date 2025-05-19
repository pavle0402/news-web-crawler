package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"crawler/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DbPool *pgxpool.Pool

const (
	maxConnections = 10
	connTimeout    = 5 * time.Second
)

func Connect(cfg *config.DBConfig) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Error parsing database URL: %v", err)
	}

	config.MaxConns = maxConnections
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	DbPool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Fatalf("Error pinging database conn: %v", err)
	}

	log.Println("Database connection established successfully.")
}
