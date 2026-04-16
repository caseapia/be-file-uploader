package image

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	group := router.Group("/storage")

	group.Post("/upload", middleware.RequirePermission(role.FileUpload), h.UploadImage)
	group.Post("/delete", middleware.RequirePermission(role.FileUpload), h.DeleteImage)
	group.Get("/list/:id", middleware.RequirePermission(role.ViewOtherFiles), h.LookupImagesByUserID)
	group.Get("/my", middleware.RequirePermission(role.ViewOwnFiles), h.LookupMyImages)
	group.Put("/album/put", middleware.RequirePermission(role.FileUpload), h.AddInAlbum)
	group.Delete("/album/delete", middleware.RequirePermission(role.FileUpload), h.RemoveFromAlbum)

	groupAdmin := router.Group("/image/admin")
	groupAdmin.Get("/list", middleware.RequirePermission(role.ManageFiles), h.LookupAllImages)
}
