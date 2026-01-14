package handlers

import (
	"context"
	"time"

	"CardFlow/internal/models"
	"CardFlow/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardHandler struct {
    service services.CardService
}

func NewCardHandler(service services.CardService) *CardHandler {
    return &CardHandler{service: service}
}

func (h *CardHandler) CreateCard(c *fiber.Ctx) error {
	var req models.CreateCardReq
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
	req.Userid = c.Locals("user_id").(uuid.UUID)

	if req.CardType == "" || req.SpendingLimit <= 0  {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "incomplete request data",
		})
	}
    res, err := h.service.CreateCard(ctx, req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "card created successfully",
		"data": res,
    })
}