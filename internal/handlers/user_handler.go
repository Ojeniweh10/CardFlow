package handlers

import (
	"CardFlow/internal/models"
	"CardFlow/internal/services"
	"CardFlow/internal/utils"
	"context"
	"time"

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
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

    if err := utils.ValidatePassword(req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
	}

    err := h.service.RegisterUser(ctx, req)
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
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

    token, err := h.service.Login(ctx, req)
    if err != nil {
        if err.Error() == "MFA required" {
            return c.Status(200).JSON(fiber.Map{
                "mfa_required": true,
                "message": "Multi-factor authentication required",
            })
        }
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "token": token,
    })
}

func (h *UserHandler) MFALogin(c *fiber.Ctx) error {
    var req models.MFALoginReq
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }

    if req.Email == "" || req.TOTPCode == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "email and TOTP code are required",
        })
    }

    token, err := h.service.MFALogin(ctx, req)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "token":   token,
        "message": "Login successful",
    })
}


func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
    user_id:= c.Locals("user_id").(uuid.UUID)
    err := h.service.VerifyEmail(ctx, user_id)
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
    var otp models.Otp
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
    if err := c.BodyParser(&otp); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
    user_id:= c.Locals("user_id").(uuid.UUID)

    if otp.Otp == ""{
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "incomplete data",
        })
    }

    err := h.service.VerifyOtp(ctx, user_id, otp.Otp)
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

func (h *UserHandler) EnableMFA(c *fiber.Ctx) error{
    user_id:= c.Locals("user_id").(uuid.UUID)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
    res, err := h.service.EnableMFA(ctx, user_id)
    if err != nil{
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": res,
    })
}

func (h *UserHandler) VerifyMFA(c *fiber.Ctx) error{
    var data models.VerifyMFA
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
    if data.TotpCode == ""{
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "incomplete data",
        })
    }
    user_id:= c.Locals("user_id").(uuid.UUID)
    err := h.service.VerifyMFA(ctx, user_id, data.TotpCode)
    if err != nil{
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "multi-factor authentication enabled",
    })
}