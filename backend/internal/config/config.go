package config

import "os"

// Config holds all application configuration.
type Config struct {
	AppEnv      string // development, production
	Port        string // server port
	DatabaseURL string // PostgreSQL connection string (compatible with Supabase)
}

// Load reads configuration from environment variables.
// Uses standard DATABASE_URL format for Supabase compatibility:
// postgres://user:password@host:port/dbname?sslmode=disable
func Load() *Config {
	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://hotel_user:hotel_password@localhost:5433/hotel_db?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
