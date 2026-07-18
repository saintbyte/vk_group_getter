package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saintbyte/vk_group_getter/internal/config"
)

func APIKey(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-API-Key")
		if key == "" || key != cfg.APIKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.Next()
	}
}
