package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSConfig returns a configured CORS middleware.
// Allows frontend development server and production origins.
func CORSConfig() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:3001, https://www.qk3houronline.com, https://qk3houronline.com, https://jw0924230.github.io",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,x-github-build-token",
		AllowCredentials: true,
	})
}
