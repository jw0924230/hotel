package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"hotel-backend/internal/database"
)

// HealthHandler holds dependencies for health check endpoints.
type HealthHandler struct {
	DB *database.DB
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{DB: db}
}

// Check returns the health status of the service and database.
// GET /api/health
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dbStatus := "connected"
	if err := h.DB.Ping(ctx); err != nil {
		dbStatus = "disconnected"
	}

	return c.JSON(fiber.Map{
		"status":    "ok",
		"service":   "hotel-backend",
		"database":  dbStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
