package album

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	group := router.Group("/album")

	group.Post("/create", middleware.RequirePermission(role.FileUpload), h.CreateAlbum)
	group.Delete("/delete/:id", middleware.RequirePermission(role.FileUpload), h.DeleteAlbum)
	group.Get("/lookup/:id", middleware.RequirePermission(role.ViewOwnFiles), h.LookupAlbum)
}
