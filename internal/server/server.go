package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/saintbyte/vk_group_getter/docs"
	"github.com/saintbyte/vk_group_getter/internal/app"
	"github.com/saintbyte/vk_group_getter/internal/handler"
	"github.com/saintbyte/vk_group_getter/internal/middleware"
	"github.com/saintbyte/vk_group_getter/internal/queue"
	swagger "github.com/swaggo/fiber-swagger"
)

func New(a *app.App, q *queue.Queue) *fiber.App {
	host := a.Config.Host
	if host == "" {
		host = "localhost"
	}
	origin := fmt.Sprintf("http://%s:%s,http://127.0.0.1:%s", host, a.Config.Port, a.Config.Port)

	fiberApp := fiber.New()
	fiberApp.Use(logger.New())
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: origin,
		AllowHeaders: "Origin, Content-Type, Accept, X-API-Key",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	h := handler.New(a, q)

	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	fiberApp.Get("/swagger/*", swagger.WrapHandler)

	fiberApi := fiberApp.Group("/", middleware.APIKey(a.Config))
	fiberApi.Post("/groups", h.AddGroup)

	return fiberApp
}
