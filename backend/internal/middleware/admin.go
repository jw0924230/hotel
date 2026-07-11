package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AdminOnly checks if the authenticated user has the 'admin' role.
// Must be placed after JWTMiddleware in the route chain.
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "權限不足，此作業僅限系統管理員執行"})
		}
		return c.Next()
	}
}
