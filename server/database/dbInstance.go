package database

import (
	"context"
	"log"
	"time"
	"sync"
	"gollama/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbInstance *pgxpool.Pool
	once       sync.Once
)

func GetDbInstance() *pgxpool.Pool {
	once.Do(initDB)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := dbInstance.Ping(ctx); err != nil {
		log.Println("Database connection lost, reinitializing pool...")
		initDB()
	}
	return dbInstance
}

func initDB() {
	connStr := config.ENV.DatabaseURL

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Unable to parse database connection string: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbInstance, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	log.Println("Database connection pool created successfully.")
}
