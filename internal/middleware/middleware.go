package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// LoggerMiddleware returns a logger middleware
func LoggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	})
}

// RecoverMiddleware returns a recover middleware
func RecoverMiddleware() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}

// Logger returns a custom logger middleware
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request details
		log.Printf("[%s] %s %s - %v - %s",
			c.Method(),
			c.Path(),
			c.IP(),
			time.Since(start),
			c.Response().Header.Peek("Content-Type"),
		)

		return err
	}
}

// Recovery returns a recovery middleware using Fiber's built-in recovery
func Recovery() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Printf("Panic recovered: %v", e)
		},
	})
}

// CORSMiddleware returns a CORS middleware
func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
	})
}

// AuthMiddleware is a sample authentication middleware
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		token := c.Get("Authorization")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization token",
			})
		}

		// Simple token validation (replace with your actual auth logic)
		if token != "Bearer valid-token" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization token",
			})
		}

		// Set user context (example)
		c.Locals("user_id", "12345")

		return c.Next()
	}
}

// RateLimitMiddleware returns a rate limiting middleware
func RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Simple rate limiting logic (you can use fiber/middleware/limiter for production)
		// This is just a placeholder implementation
		return c.Next()
	}
}
