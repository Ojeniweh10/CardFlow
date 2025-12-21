package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type loginAttempt struct {
    Count     int
    ExpiresAt time.Time
}


var loginAttempts = make(map[string]*loginAttempt)
var mu sync.Mutex



func LoginRateLimit() fiber.Handler {
    return func(c *fiber.Ctx) error {
        ip := c.IP()
        now := time.Now()

        mu.Lock()
        defer mu.Unlock()

        attempt, exists := loginAttempts[ip]

        // First attempt OR expired window
        if !exists || now.After(attempt.ExpiresAt) {
            loginAttempts[ip] = &loginAttempt{
                Count:     1,
                ExpiresAt: now.Add(15 * time.Minute),
            }
            return c.Next()
        }

        // Too many attempts
        if attempt.Count >= 5 {
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "success": false,
                "message": "Too many login attempts. Please try again later.",
            })
        }

        // Increment attempt count
        attempt.Count++
        return c.Next()
    }
}
