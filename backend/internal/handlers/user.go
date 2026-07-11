package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"hotel-backend/internal/database"
	"hotel-backend/internal/models"
)

// UserHandler handles HTTP requests for user management.
type UserHandler struct {
	DB *database.DB
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// UserCreateRequest defines fields to create a user.
type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserUpdateRequest defines fields to update a user.
type UserUpdateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // Optional password change
	Role     string `json:"role"`
}

// List lists all users.
// GET /api/users
func (h *UserHandler) List(c *fiber.Ctx) error {
	ctx := context.Background()
	query := `SELECT id, email, role, created_at FROM users ORDER BY id ASC`

	rows, err := h.DB.Pool.Query(ctx, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to query users: %v", err)})
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to scan user: %v", err)})
		}
		users = append(users, user)
	}

	return c.JSON(users)
}

// Create creates a new admin or vendor user.
// POST /api/users
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Role = strings.TrimSpace(strings.ToLower(req.Role))

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and Password are required"})
	}

	if req.Role != "admin" && req.Role != "vendor" {
		req.Role = "vendor" // Default role
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}

	ctx := context.Background()
	var user models.User
	query := `
		INSERT INTO users (email, password, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, role, created_at`

	err = h.DB.Pool.QueryRow(ctx, query, req.Email, string(hashedPassword), req.Role).Scan(
		&user.ID, &user.Email, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "此 Email 帳號已被使用"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to create user: %v", err)})
	}

	return c.JSON(user)
}

// Update updates an existing user's details.
// PUT /api/users/:id
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	var req UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Role = strings.TrimSpace(strings.ToLower(req.Role))

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	if req.Role != "admin" && req.Role != "vendor" {
		req.Role = "vendor"
	}

	ctx := context.Background()
	var err error

	if req.Password != "" {
		// Update details with password change
		hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if errHash != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
		}

		query := `UPDATE users SET email = $1, password = $2, role = $3 WHERE id = $4`
		_, err = h.DB.Pool.Exec(ctx, query, req.Email, string(hashedPassword), req.Role, id)
	} else {
		// Update details only (no password change)
		query := `UPDATE users SET email = $1, role = $2 WHERE id = $3`
		_, err = h.DB.Pool.Exec(ctx, query, req.Email, req.Role, id)
	}

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "此 Email 帳號已被使用"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to update user: %v", err)})
	}

	return c.JSON(fiber.Map{"message": "使用者帳號更新成功"})
}

// Delete deletes a user.
// DELETE /api/users/:id
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	// Prevent admin from deleting themselves
	currentUserId := fmt.Sprintf("%v", c.Locals("userId"))
	if id == currentUserId {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "您無法刪除自己目前登入的帳號"})
	}

	ctx := context.Background()
	query := `DELETE FROM users WHERE id = $1`
	_, err := h.DB.Pool.Exec(ctx, query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("failed to delete user: %v", err)})
	}

	return c.JSON(fiber.Map{"message": "使用者帳號已成功刪除"})
}
