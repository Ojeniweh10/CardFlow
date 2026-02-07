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

func (h *CardHandler)FetchAllCards(c *fiber.Ctx) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Userid := c.Locals("user_id").(uuid.UUID)
	res, err := h.service.GetAllCards(ctx, Userid)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "cards fetched successfully",
		"data": res,
    })
}

func (h *CardHandler)FetchCardById(c *fiber.Ctx) error{
    var req models.GetCardReq
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req.UserId = c.Locals("user_id").(uuid.UUID)
    req.CardId = c.Params("id")
	res, err := h.service.GetCardById(ctx, req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "card fetched successfully",
		"data": res,
    })
}

func (h *CardHandler)ModifyCardStatus(c *fiber.Ctx) error{
    var req models.GetCardReq
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
    Status := c.Params("status")
    if Status == "" || req.CardId == ""{
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "incomplete request data",
		})
    }
    err := h.service.ModifyCardStatus(ctx, req, Status)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "status modified successfully",
		"data": nil,
    })
}

func (h *CardHandler)TopUpCard(c *fiber.Ctx) error{
    var req models.TopUpCardReq
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
    req.Userid = c.Locals("user_id").(uuid.UUID)
    req.Cardid = c.Params("id")
    if req.Amount <= 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "incomplete request data",
		})
    }
    res, err := h.service.TopUpCard(ctx, req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "card funded successfully",
		"data": res,
    })
}
