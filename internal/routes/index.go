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
    api.Post("/login",middleware.LoginRateLimit(), userHandler.Login)//login
    api.Post("/login/mfa",middleware.LoginRateLimit(), userHandler.MFALogin)//login with mfa
    api.Post("/verify",middleware.JWTProtected(),userHandler.VerifyEmail)//verify email otp will be sent
    api.Post("/otp", middleware.JWTProtected(), userHandler.VerifyOtp)//verify otp which is gotten from email
    api.Post("/mfa/setup", middleware.JWTProtected(), userHandler.EnableMFA)//this for user to set up mfa
    api.Post("/mfa/verify", middleware.JWTProtected(), userHandler.VerifyMFA)//this is for user to verify mfa totp 
    api.Patch("/change-password", middleware.JWTProtected(), userHandler.ChangePassword)//change password
    api.Put("/reset-password", userHandler.ResetPassword)//password is then reset and you can log in
    api.Put("/forgot-pwd/otp", userHandler.ForgotPwdOtp)//forgot password otp is verified
    api.Post("/forgot-password",userHandler.ForgotPassword)//forgot password while trying to log in, otp is sent to email
    api.Post("/", userHandler.CreateUser)//create user
    
}

func KycRoutes(app *fiber.App, db *gorm.DB){
    userRepo := repositories.NewUserRepository(db)
    kycRepo := repositories.NewKycRepository(db)
    kycService := services.NewKycService(kycRepo, userRepo)
    kycHandler := handlers.NewKycHandler(kycService)
    api := app.Group("/api/v1/kyc")
    api.Post("/selfie", middleware.JWTProtected(), kycHandler.Uploadimage)//endpoint to upload selfie image
    api.Post("/document", middleware.JWTProtected(), kycHandler.UploadKycDocument)//a picture of an id, nin or voters card or passport
    api.Post("/proof-of-address", middleware.JWTProtected(), kycHandler.UploadProofOfAddress)//upload proof of address doc, kyc is marked as verified after this
}

func CardRoutes(app *fiber.App, db *gorm.DB) {
    cardRepo := repositories.NewCardRepository(db)
    kycRepo := repositories.NewKycRepository(db)
    userRepo:= repositories.NewUserRepository(db)
    txnRepo := repositories.NewTransactionRepository(db)
    cardService := services.NewCardService(userRepo, kycRepo, cardRepo, txnRepo)
    cardHandler := handlers.NewCardHandler(cardService)

    api := app.Group("/api/v1/cards")
    api.Post("/top-up/:id", middleware.JWTProtected(), cardHandler.TopUpCard)//fund card
    api.Patch("/:status", middleware.JWTProtected(), cardHandler.ModifyCardStatus)//modify card status, freeze, delete etc
    api.Get("/:id",middleware.JWTProtected(), cardHandler.FetchCardById)//gets all details of a particular card for a user
    api.Get("/", middleware.JWTProtected(), cardHandler.FetchAllCards)//gets users card(s)
    api.Post("/",middleware.JWTProtected(), cardHandler.CreateCard)//create card
}

func TransactionRoutes(app *fiber.App, db *gorm.DB) {
    cardRepo := repositories.NewCardRepository(db)
    userRepo := repositories.NewUserRepository(db)
    transactionRepo := repositories.NewTransactionRepository(db)
    transactionService := services.NewTransactionService(transactionRepo, cardRepo, userRepo)
    transactionHandler := handlers.NewTransactionHandler(transactionService)

    api := app.Group("/api/v1/transactions")
    api.Post("/webhook", transactionHandler.HandleWebhook)// receive authorize and capture, this is practically going to simulate an incoming card transactions
    api.Get("/:id",middleware.JWTProtected(), transactionHandler.GetCardTransactions)//get all transactions for a particular card
}
