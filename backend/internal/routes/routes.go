package routes

import (
	"github.com/gofiber/fiber/v2"

	"hotel-backend/internal/database"
	"hotel-backend/internal/handlers"
	"hotel-backend/internal/middleware"
)

// Setup registers all application routes.
func Setup(app *fiber.App, db *database.DB) {
	// API group
	api := app.Group("/api")

	// Health check
	healthHandler := handlers.NewHealthHandler(db)
	api.Get("/health", healthHandler.Check)

	// Auth routes
	authHandler := handlers.NewAuthHandler(db)
	api.Post("/auth/login", authHandler.Login)

	// Hotel CRUD routes
	hotelHandler := handlers.NewHotelHandler(db)
	api.Get("/hotels", hotelHandler.List)
	api.Get("/hotels/:id", hotelHandler.Get)

	// Homepage featured hotels routes
	homepageHandler := handlers.NewHomepageHotelHandler(db)
	api.Get("/homepage-hotels", homepageHandler.GetHomepageHotels)
	api.Put("/homepage-hotels", middleware.JWTMiddleware(), homepageHandler.UpdateHomepageHotels)

	// Categories routes (replaces old /api/cities)
	categoryHandler := handlers.NewCategoryHandler(db)
	api.Get("/categories", categoryHandler.List)
	api.Get("/regions", categoryHandler.Regions)
	api.Get("/regions/combined", categoryHandler.CombinedRegions)
	api.Post("/categories", middleware.JWTMiddleware(), categoryHandler.Create)
	api.Put("/categories/:id", middleware.JWTMiddleware(), categoryHandler.Update)
	api.Delete("/categories/:id", middleware.JWTMiddleware(), categoryHandler.Delete)

	// Protected routes (require JWT)
	api.Put("/hotels/:id", middleware.JWTMiddleware(), hotelHandler.Upsert)

	// Blog post routes
	postHandler := handlers.NewPostHandler(db)
	api.Get("/posts", postHandler.List)
	api.Get("/posts/:id", postHandler.Get)
	api.Post("/posts", middleware.JWTMiddleware(), postHandler.Create)
	api.Put("/posts/:id", middleware.JWTMiddleware(), postHandler.Update)
	api.Delete("/posts/:id", middleware.JWTMiddleware(), postHandler.Delete)

	// User account management routes (admin only)
	userHandler := handlers.NewUserHandler(db)
	api.Get("/users", middleware.JWTMiddleware(), middleware.AdminOnly(), userHandler.List)
	api.Post("/users", middleware.JWTMiddleware(), middleware.AdminOnly(), userHandler.Create)
	api.Put("/users/:id", middleware.JWTMiddleware(), middleware.AdminOnly(), userHandler.Update)
	api.Delete("/users/:id", middleware.JWTMiddleware(), middleware.AdminOnly(), userHandler.Delete)

	// Frontend deploy trigger route (admin only)
	deployHandler := handlers.NewDeployHandler()
	api.Post("/deploy", middleware.JWTMiddleware(), middleware.AdminOnly(), deployHandler.TriggerDeploy)
}


