package handlers

import (
	"CardFlow/internal/models"
	"CardFlow/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
    service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
    return &UserHandler{service: service}
}


func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req models.CreateUserRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }

	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "incomplete request data",
		})
	}

    err := h.service.RegisterUser(req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
		"message": "user created successfully",
    })
}


func (h *UserHandler) Login(c *fiber.Ctx) error {
    var req models.LoginReq
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }

    if req.Email == "" || req.Password == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "email and password are required",
        })
    }

    token, err := h.service.Login(req)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "token": token,
    })
}

func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
    user_id:= c.Locals("user_id").(uuid.UUID)
    err := h.service.VerifyEmail(user_id)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "otp sent successfully",
    })
}

func (h *UserHandler) VerifyOtp(c *fiber.Ctx) error {
    var otp string
    if err := c.BodyParser(&otp); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
    user_id:= c.Locals("user_id").(uuid.UUID)

    if otp == ""{
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "incomplete data",
        })
    }

    err := h.service.VerifyOtp(user_id, otp)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "email verified successfully",
    })

}