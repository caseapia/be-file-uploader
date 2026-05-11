package user

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	group := router.Group("/user")
	groupAdmin := group.Group("/admin", middleware.Require([]string{"not_banned"}))

	group.Get("/me", h.LookupMyAccount)
	group.Get("/lookup/:id",
		middleware.RequirePermission(role.ViewOtherProfiles),
		middleware.Require([]string{"not_banned"}),
		h.LookupProfile,
	)
	group.Get("/shareX/generate",
		middleware.RequirePermission(role.FileUpload),
		middleware.Require([]string{"not_banned"}),
		middleware.Require([]string{"not_banned", "not_restricted_upload"}),
		h.GenerateAPIToken,
	)
	group.Get("/lookupByName/:name", h.LookupUsersByPart)
	group.Patch("/avatar",
		middleware.Require([]string{"not_banned", "not_restricted_upload"}),
		h.UploadAvatar,
	)

	groupAdmin.Get("/users", middleware.RequirePermission(role.ManageUsers), h.PopulateUserList)
	groupAdmin.Put("/role/add", middleware.RequirePermission(role.ManageUsers), h.AddUserInRole)
	groupAdmin.Delete("/role/delete", middleware.RequirePermission(role.ManageUsers), h.RemoveUserFromRole)
	groupAdmin.Patch("/storage-limit/update", middleware.RequirePermission(role.ManageUsers), h.SetUploadLimit)
	groupAdmin.Patch("/verify/:id", middleware.RequirePermission(role.ManageUsers), h.VerifyUser)
	groupAdmin.Delete("/shareX/reset/:id", middleware.RequirePermission(role.ManageUsers), h.ResetUserAPIToken)
	groupAdmin.Post("/restriction", middleware.RequirePermission(role.ManageUsers), h.RestrictUser)
	groupAdmin.Get("/restrictions/list/:id", middleware.RequirePermission(role.ManageUsers), h.LookupUserRestrictions)
	groupAdmin.Delete("/restriction/delete/:id", middleware.RequirePermission(role.ManageUsers), h.RemoveUserRestriction)
}
