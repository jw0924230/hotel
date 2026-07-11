package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://hotel_user:hotel_password@localhost:5433/hotel_db?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Query indexes on hotel_images table
	rows, err := pool.Query(ctx, `
		SELECT indexname, indexdef 
		FROM pg_indexes 
		WHERE tablename = 'hotel_images'`)
	if err != nil {
		log.Fatalf("Failed to query indexes: %v", err)
	}
	defer rows.Close()

	log.Println("--- Indexes on hotel_images ---")
	for rows.Next() {
		var name, def string
		_ = rows.Scan(&name, &def)
		log.Printf("Index: %s -> %s", name, def)
	}
}
