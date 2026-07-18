package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/saintbyte/vk_group_getter/internal/app"
	"github.com/saintbyte/vk_group_getter/internal/model"
	"github.com/saintbyte/vk_group_getter/internal/queue"
)

type Handler struct {
	App   *app.App
	Queue *queue.Queue
}

func New(a *app.App, q *queue.Queue) *Handler {
	return &Handler{App: a, Queue: q}
}

func normalizeAddress(addr string) (string, string) {
	addr = strings.TrimSpace(addr)
	addr = strings.TrimRight(addr, "/")
	addr_type := "vk_group"
	lower := strings.ToLower(addr)
	vkPosiblePrefix := []string{
		"https://vk.com/",
		"http://vk.com/",
		"vk.com/",
		"m.vk.com/",
		"https://m.vk.com/",
		"http://m.vk.com/",
		"https://vk.ru/",
		"http://vk.ru/",
		"vk.ru/",
		"m.vk.ru/",
		"https://m.vk.ru/",
		"http://m.vk.ru/",
	}
	for _, prefix := range vkPosiblePrefix {
		if strings.HasPrefix(lower, prefix) {
			addr = addr[len(prefix):]
			break
		}
	}
	addr = strings.TrimLeft(addr, "/")
	return addr, addr_type
}

// AddGroup godoc
// @Summary      Add a VK group
// @Description  Add a VK group address for processing
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        X-API-Key header string true "API key"
// @Param        address body object true "VK group address"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /groups [post]
func (h *Handler) AddGroup(c *fiber.Ctx) error {
	var req struct {
		Address string `json:"address"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.Address == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "address is required"})
	}

	address, tp := normalizeAddress(req.Address)
	if address == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid address"})
	}

	group := model.Group{
		Address: address,
		Type:    tp,
	}
	if err := h.App.DB.Create(&group).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save group"})
	}

	h.Queue.Push(group.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": group.ID, "address": group.Address})
}
