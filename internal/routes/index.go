package routes

import (
	"CardFlow/internal/handlers"
	"CardFlow/internal/middleware"
	"CardFlow/internal/repositories"
	"CardFlow/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
    UserRoutes(app, db)
    //CardRoutes(app, db)
    //TransactionRoutes(app, db)
}



func UserRoutes(app *fiber.App, db *gorm.DB) {
    userRepo := repositories.NewUserRepository(db)
    userService := services.NewUserService(userRepo)
    userHandler := handlers.NewUserHandler(userService)

    api := app.Group("/api/v1/users")
    api.Post("/", userHandler.CreateUser)
    api.Post("/login",middleware.LoginRateLimit(), userHandler.Login)
    api.Post("/login/mfa",middleware.LoginRateLimit(), userHandler.MFALogin)
    api.Post("/verify",middleware.JWTProtected(),userHandler.VerifyEmail)
    api.Post("/otp", middleware.JWTProtected(), userHandler.VerifyOtp)
    api.Post("/mfa/setup", middleware.JWTProtected(), userHandler.EnableMFA)
    api.Post("/mfa/verify", middleware.JWTProtected(), userHandler.VerifyMFA)
}

// func CardRoutes(app *fiber.App, db *gorm.DB) {
//     cardRepo := repositories.NewCardRepository(db)
//     cardHandler := handlers.NewCardHandler(cardRepo)

//     api := app.Group("/api/v1/cards")
//     api.Post("/", cardHandler.CreateCard)
//     api.Get("/:id", cardHandler.GetCardById)
// }

// func TransactionRoutes(app *fiber.App, db *gorm.DB) {
//     transactionRepo := repositories.NewTransactionRepository(db)
//     transactionHandler := handlers.NewTransactionHandler(transactionRepo)

//     api := app.Group("/api/v1/transactions")
//     api.Post("/", transactionHandler.CreateTransaction)
//     api.Get("/:id", transactionHandler.GetTransactionById)
// }
