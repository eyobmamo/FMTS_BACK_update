// package config

// import (
// 	"context"
// 	"log"
// 	"os"

// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/joho/godotenv"
// )

// func LoadConfig() string {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	return os.Getenv("postgressURL")
// }

// func ConnectSupabasePool(connString string) *pgxpool.Pool {
// 	ctx := context.Background()

// 	cfg, err := pgxpool.ParseConfig(connString)
// 	if err != nil {
// 		log.Fatalf("Unable to parse database config: %v\n", err)
// 	}

// 	// Use simple protocol to avoid server-side prepared statements with PgBouncer (transaction pooling)
// 	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

// 	pool, err := pgxpool.NewWithConfig(ctx, cfg)
// 	if err != nil {
// 		log.Fatalf("Unable to connect to database pool: %v\n", err)
// 	}
// 	return pool
// }

package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func LoadConfig() string {
	// Try to load .env (for local dev only)
	_ = godotenv.Load() // ignore if missing

	conn := os.Getenv("postgressURL")
	if conn == "" {
		log.Fatal("Database connection string not found. Please set 'postgressURL' environment variable.")
	}
	return conn
}

func ConnectSupabasePool(connString string) *pgxpool.Pool {
	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v\n", err)
	}

	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Unable to connect to database pool: %v\n", err)
	}
	return pool
}
