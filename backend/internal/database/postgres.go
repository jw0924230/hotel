package database

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schemaSQL string

//go:embed townships.json
var townshipsJSON []byte

type townshipSeed struct {
	City      string `json:"city"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	SortOrder int    `json:"sort_order"`
}

// DB wraps pgxpool.Pool to provide database operations.
type DB struct {
	Pool *pgxpool.Pool
}

// Connect creates a new connection pool to PostgreSQL.
// The connString should be a standard PostgreSQL connection URL:
// postgres://user:password@host:port/dbname?sslmode=disable
//
// This format is directly compatible with Supabase connection strings.
func Connect(ctx context.Context, connString string) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Connection pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// InitSchema executes the embedded SQL schema to create tables and indexes.
func (db *DB) InitSchema(ctx context.Context) error {
	if _, err := db.Pool.Exec(ctx, schemaSQL); err != nil {
		return fmt.Errorf("failed to run database init SQL: %w", err)
	}

	return db.seedTownships(ctx)
}

func (db *DB) seedTownships(ctx context.Context) error {
	var seeds []townshipSeed
	if err := json.Unmarshal(townshipsJSON, &seeds); err != nil {
		return fmt.Errorf("failed to decode township seed data: %w", err)
	}
	if len(seeds) != 368 {
		return fmt.Errorf("invalid township seed count: got %d, want 368", len(seeds))
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin township seed transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, seed := range seeds {
		result, err := tx.Exec(ctx, `
			INSERT INTO categories (type, name, parent_id, external_code, sort_order)
			SELECT 'township', $1, c.id, $2, $3
			FROM categories c
			WHERE c.type = 'city' AND c.name = $4 AND c.parent_id IS NULL
			ON CONFLICT (external_code) WHERE external_code IS NOT NULL DO UPDATE SET
				name = EXCLUDED.name,
				parent_id = EXCLUDED.parent_id,
				sort_order = EXCLUDED.sort_order`,
			seed.Name, seed.Code, seed.SortOrder, seed.City,
		)
		if err != nil {
			return fmt.Errorf("failed to seed township %s %s: %w", seed.City, seed.Name, err)
		}
		if result.RowsAffected() == 0 {
			return fmt.Errorf("failed to find city category for township %s %s", seed.City, seed.Name)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit township seeds: %w", err)
	}
	return nil
}

// Close closes the connection pool.
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Ping checks if the database connection is alive.
func (db *DB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}
