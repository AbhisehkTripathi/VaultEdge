package main

import (
	"github.com/gofiber/fiber/v2"

	"UploadDocument-Saas/internal/middleware"
	"UploadDocument-Saas/pkg/logger"
	"UploadDocument-Saas/routes"
)

func main() {
	logFile := logger.InitLogFile()
	defer logFile.Close()

	app := fiber.New()

	// Middleware
	app.Use(middleware.RecoverMiddleware())
	app.Use(middleware.LoggerMiddleware())

	// Load routes
	routes.SetupRoutes(app)

	app.Listen(":3000")
}
