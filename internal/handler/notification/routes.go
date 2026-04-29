package notification

import "github.com/gofiber/fiber/v3"

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	group := router.Group("/notifications")

	group.Get("/my", h.SearchMyNotifications)
	group.Patch("/action/read/:id", h.ReadNotification)
}
