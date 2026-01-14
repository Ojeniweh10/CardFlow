package main

import (
	"CardFlow/internal/config"
	database "CardFlow/internal/database"
	"CardFlow/internal/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// 1. Connect to database using GORM
	db := database.NewGormConnection()

	// 2. Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "CardFlow Service",
	})

	// 3. Logger middleware
	app.Use(fiberlogger.New())

	// 4. Health check
	app.Get("/admin/healthchecker", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Welcome to CardFlow",
		})
	})

	// 5. Gateway auth middleware
	app.Use(func(c *fiber.Ctx) error {
		auth := c.Get("G-Auth")
		if auth == "" || auth != config.GatewaySecret {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Access denied",
			})
		}
		return c.Next()
	})

	// 6. Route registration (dependency injection)
	routes.Routes(app, db)

	// 7. 404 handler
	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Route Not Found",
		})
	})

	// 8. Graceful shutdown
	go func() {
		if err := app.Listen(":8081"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server stopped: %v", err)
		}
	}()

	log.Println("Server started on :8081")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // CTRL+C or kill signal
	<-quit

	log.Println("Shutting down gracefully...")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(timeoutCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
