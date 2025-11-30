package routes

import (
	"CardFlow/internal/handlers"
	"CardFlow/internal/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {

	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	api := app.Group("/api/v1/users")

	api.Post("/", userHandler.CreateUser)
	api.Get("/:id", userHandler.GetUserById)
}