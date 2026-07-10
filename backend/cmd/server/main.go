package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"hotel-backend/internal/config"
	"hotel-backend/internal/database"
	appMiddleware "hotel-backend/internal/middleware"
	"hotel-backend/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Connected to PostgreSQL")

	// Automatically run database initialization SQL
	if err := db.InitSchema(ctx); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
	log.Println("✅ Database schema initialized successfully")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Hotel Backend",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	// Global middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(appMiddleware.CORSConfig())

	// Register routes
	routes.Setup(app, db)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("🚀 Server started on port %s (env: %s)", cfg.Port, cfg.AppEnv)

	<-quit
	log.Println("Shutting down server...")

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
