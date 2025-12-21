package handlers

import (
	"CardFlow/internal/models"
	"CardFlow/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type KycHandler struct {
    service services.KycService
}

func NewKycHandler(service services.KycService) *KycHandler {
    return &KycHandler{service: service}
}

func (h *KycHandler) Uploadimage(c *fiber.Ctx) error{
	var data models.KycProfile
	if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
	data.Userid = c.Locals("user_id").(uuid.UUID)

	if data.DOB == "" || data.ImageStr == ""{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "incomplete data",
        })
	}

	err := h.service.Uploadimage(data)
	if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "image processed successfully",
    })
}

func (h *KycHandler)VerifyBVN(c *fiber.Ctx) error{
    var data models.Bvn
	if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
	data.Userid = c.Locals("user_id").(uuid.UUID)
    if data.Bvn == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "incomplete data",
        })
	}
    err := h.service.VerifyBVN(data)
	if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "bvn verified successfully",
    })
}