package roadmap

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	group := router.Group("/roadmap")

	group.Get("/list", h.GetRoadmap)
}

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	admin := router.Group("/roadmap/admin")

	admin.Post("/task/add", middleware.RequirePermission(role.Developer), h.AddRoadmapTask)
	admin.Patch("/task/edit", middleware.RequirePermission(role.Developer), h.EditRoadmapTask)
}
