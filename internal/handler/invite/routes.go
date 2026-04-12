package invite

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	groupAdmin := router.Group("/invite/admin")

	groupAdmin.Get("/list", middleware.RequirePermission(role.ManageUsers), h.SearchAllInvites)
	groupAdmin.Post("/create", middleware.RequirePermission(role.ManageUsers), h.CreateInvite)
	groupAdmin.Delete("/revoke", middleware.RequirePermission(role.ManageUsers), h.RevokeInvite)
}
