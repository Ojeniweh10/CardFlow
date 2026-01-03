package handlers

import (
	"CardFlow/internal/models"
	"CardFlow/internal/services"
	"context"
	"time"

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
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

	err := h.service.Uploadimage(ctx, data)
	if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "image uploaded successfully",
    })
}

func (h *KycHandler)UploadKycDocument(c *fiber.Ctx) error{
    var data models.KycDoc
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
	data.Userid = c.Locals("user_id").(uuid.UUID)
    if data.DocStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "incomplete data",
        })
	}
    err := h.service.UploadKycDocument(ctx, data)
	if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "Document uploaded successfully",
    })
}

func (h *KycHandler)UploadProofOfAddress(c *fiber.Ctx) error{
    var data models.KycDoc
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
	data.Userid = c.Locals("user_id").(uuid.UUID)
    if data.DocStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "incomplete data",
        })
	}
    err := h.service.UploadProofOfAddress(ctx, data)
	if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "Proof of Address uploaded successfully, Verification Pending",
    })
}