package user

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	group := router.Group("/user")

	group.Get("/me", h.LookupMyAccount)
	group.Get("/lookup/:id", h.LookupProfile)
}
