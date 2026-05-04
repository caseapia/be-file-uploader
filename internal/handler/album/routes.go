package album

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	group := router.Group("/album")
	action := group.Group("/action")

	action.Post("/create", middleware.RequirePermission(role.FileUpload), h.CreateAlbum)
	action.Delete("/delete/:id", middleware.RequirePermission(role.FileUpload), h.DeleteAlbum)
	group.Get("/lookup/:id", middleware.RequirePermission(role.ViewOwnFiles), h.LookupAlbum)
	group.Get("/lookupAll", middleware.RequirePermission(role.ManageFiles), h.AllAlbums)
}
