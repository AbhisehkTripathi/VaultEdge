package routes

import (
	"github.com/gofiber/fiber/v2"

	"UploadDocument-Saas/internal/handlers"
	"UploadDocument-Saas/internal/middleware"
	"UploadDocument-Saas/internal/websocket"
)

// SetupRoutes registers all application routes.
func SetupRoutes(app *fiber.App) {
	// Global Middleware
	app.Use(middleware.CORSMiddleware()) // CORS support
	app.Use(middleware.Logger())         // Custom logger
	app.Use(middleware.Recovery())       // Recover from panics

	// Health Check (public endpoint)
	app.Get("/health", handlers.HealthCheck)

	// WebSocket for real-time communication
	app.Get("/ws", websocket.HandleWebSocket)

	// Static file serving for uploads
	app.Static("/uploads", "./uploads")

	// API routes group with rate limiting
	api := app.Group("/api")
	api.Use(middleware.RateLimitMiddleware())

	// Document routes
	document := api.Group("/document")
	document.Post("/upload", handlers.UploadDocument)
	document.Get("/:id", handlers.GetDocumentByID)
	document.Get("/", handlers.ListDocuments)

	// Folder routes
	folder := api.Group("/folder")
	folder.Get("/", handlers.ListFolders)

	// Master routes
	master := api.Group("/master")
	master.Get("/", handlers.ListMasters)

	// Protected routes (require authentication)
	protected := api.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	protected.Get("/profile", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{
			"message": "This is a protected route",
			"user_id": userID,
		})
	})

	// Kafka testing (optional)
	app.Post("/kafka/send", handlers.SendKafkaTestMessage)
}
