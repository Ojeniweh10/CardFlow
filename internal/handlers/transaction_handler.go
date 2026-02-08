package handlers

import (
	"CardFlow/internal/models"
	"CardFlow/internal/services"
	"CardFlow/internal/utils"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionHandler struct {
    service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
    return &TransactionHandler{service: service}
}



func(h *TransactionHandler)HandleWebhook(c *fiber.Ctx) error{
	var data models.WebhookReq
	rawBody := c.Body()
    if len(rawBody) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "empty request body",
        })
    }
    fmt.Println("Received body:")
    fmt.Printf("%q\n", string(rawBody))
    fmt.Println("Received signature:", c.Get("X-Signature"))
    signature := c.Get("X-Signature")
    if signature == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "missing hmac signature",
        })
    }
    if err := utils.ValidateHMAC(rawBody, signature); err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "invalid hmac signature",
        })
    }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "invalid request body",
        })
    }
    if data.Amount <=0 || data.CardReference == "" || data.Currency == "" || data.Direction == "" || data.IdempotencyKey == "" || data.Status == "" || data.TransactionID == "" || data.Type == "" || data.Timestamp.IsZero() {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "incomplete request data",
        })
    }
	res, err := h.service.WebhookTransaction(ctx, data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "data processed successfully",
		"data": res,
    })
}

func (h *TransactionHandler)GetCardTransactions(c *fiber.Ctx) error{
    var data models.GetCardTransactionsReq
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
    data.Userid = c.Locals("user_id").(uuid.UUID)
    data.Cardid = c.Params("id")
    if data.Cardid == ""{
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "incomplete request data",
        })
    }
    res, err := h.service.GetCardTransactions(ctx, data)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "message": "data processed successfully",
		"data": res,
    })
}