package auth

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	group := router.Group("/auth")

	group.Post("/register", h.Register)
	group.Post("/login", h.Login)
	group.Post("/refresh", h.Refresh)
}

func (h *Handler) RegisterPrivateRoutes(_ fiber.Router) {

}
