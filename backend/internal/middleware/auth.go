package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware validates the bearer token in the Authorization header.
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "未提供驗證金鑰"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "驗證格式錯誤"})
		}

		tokenString := parts[1]

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "hotel-backend-secret-key-123456" // Default fallback
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "無效的驗證金鑰"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "無效的 Claims 格式"})
		}

		// Store claims in context locals
		c.Locals("userId", claims["sub"])
		c.Locals("email", claims["email"])
		c.Locals("role", claims["role"])

		// Generate refreshed token with sliding 24h expiration
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   claims["sub"],
			"email": claims["email"],
			"role":  claims["role"],
			"exp":   time.Now().Add(24 * time.Hour).Unix(),
		})
		newTokenString, err := newToken.SignedString([]byte(jwtSecret))
		if err == nil {
			c.Set("X-Refresh-Token", newTokenString)
		}

		return c.Next()
	}
}
