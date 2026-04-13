package role

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	groupAdmin := router.Group("/roles/admin")

	groupAdmin.Get("/all", middleware.RequirePermission(role.ManageRoles), h.LookupAllRoles)
	groupAdmin.Post("/create", middleware.RequirePermission(role.ManageRoles), h.CreateRole)
	groupAdmin.Patch("/edit", middleware.RequirePermission(role.ManageRoles), h.UpdateRole)
	groupAdmin.Delete("/delete", middleware.RequirePermission(role.ManageRoles), h.DeleteRole)
}
