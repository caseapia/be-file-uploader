package user

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	group := router.Group("/user")
	groupAdmin := group.Group("/admin")

	group.Get("/me", h.LookupMyAccount)
	group.Get("/lookup/:id", h.LookupProfile)

	groupAdmin.Get("/users", middleware.RequirePermission(role.ManageUsers), h.PopulateUserList)
	groupAdmin.Put("/role/add", middleware.RequirePermission(role.ManageUsers), h.AddUserInRole)
	groupAdmin.Delete("/role/delete", middleware.RequirePermission(role.ManageUsers), h.RemoveUserFromRole)
	groupAdmin.Patch("/storage-limit/update", middleware.RequirePermission(role.ManageUsers), h.SetUploadLimit)
	groupAdmin.Patch("/verify/:id", middleware.RequirePermission(role.ManageUsers), h.VerifyUser)
}
