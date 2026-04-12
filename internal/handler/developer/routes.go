package developer

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Get("/ping", h.Ping)
}
