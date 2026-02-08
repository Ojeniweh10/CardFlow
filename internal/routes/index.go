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
    KycRoutes(app, db)
    CardRoutes(app, db)
    TransactionRoutes(app, db)
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

func KycRoutes(app *fiber.App, db *gorm.DB){
    userRepo := repositories.NewUserRepository(db)
    kycRepo := repositories.NewKycRepository(db)
    kycService := services.NewKycService(kycRepo, userRepo)
    kycHandler := handlers.NewKycHandler(kycService)
    api := app.Group("/api/v1/kyc")
    api.Post("/selfie", middleware.JWTProtected(), kycHandler.Uploadimage)
    api.Post("/document", middleware.JWTProtected(), kycHandler.UploadKycDocument)//a picture of an id, nin or voters card or passport
    api.Post("/proof-of-address", middleware.JWTProtected(), kycHandler.UploadProofOfAddress)
}

func CardRoutes(app *fiber.App, db *gorm.DB) {
    cardRepo := repositories.NewCardRepository(db)
    kycRepo := repositories.NewKycRepository(db)
    userRepo:= repositories.NewUserRepository(db)
    txnRepo := repositories.NewTransactionRepository(db)
    cardService := services.NewCardService(userRepo, kycRepo, cardRepo, txnRepo)
    cardHandler := handlers.NewCardHandler(cardService)

    api := app.Group("/api/v1/cards")
    api.Post("/top-up/:id", middleware.JWTProtected(), cardHandler.TopUpCard)
    api.Patch("/:status", middleware.JWTProtected(), cardHandler.ModifyCardStatus)
    api.Get("/:id",middleware.JWTProtected(), cardHandler.FetchCardById)
    api.Get("/", middleware.JWTProtected(), cardHandler.FetchAllCards)
    api.Post("/",middleware.JWTProtected(), cardHandler.CreateCard)
}

func TransactionRoutes(app *fiber.App, db *gorm.DB) {
    cardRepo := repositories.NewCardRepository(db)
    userRepo := repositories.NewUserRepository(db)
    transactionRepo := repositories.NewTransactionRepository(db)
    transactionService := services.NewTransactionService(transactionRepo, cardRepo, userRepo)
    transactionHandler := handlers.NewTransactionHandler(transactionService)

    api := app.Group("/api/v1/transactions")
    api.Post("/webhook", transactionHandler.HandleWebhook)// receive authorize and capture
    api.Get("/:id",middleware.JWTProtected(), transactionHandler.GetCardTransactions)
}
