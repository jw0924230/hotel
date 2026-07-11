package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"hotel-backend/internal/database"
	"hotel-backend/internal/models"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	DB *database.DB
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// Login verifies credentials and generates a JWT.
// POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	ctx := context.Background()
	var user models.User
	var passwordHash string

	query := `SELECT id, email, password, role, created_at FROM users WHERE email = $1`
	err := h.DB.Pool.QueryRow(ctx, query, req.Email).Scan(
		&user.ID, &user.Email, &passwordHash, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "帳號或密碼錯誤"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("database error: %v", err)})
	}

	// Compare bcrypt password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "帳號或密碼錯誤"})
	}

	// Generate JWT Token
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "hotel-backend-secret-key-123456" // Default fallback
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to generate token: %v", err)})
	}

	return c.JSON(models.LoginResponse{
		Token: tokenString,
		User:  user,
	})
}
