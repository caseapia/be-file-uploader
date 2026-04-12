package image

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	group := router.Group("/storage")

	group.Post("/upload", h.UploadImage)
	group.Post("/delete", h.DeleteImage)
	group.Get("/my", middleware.RequirePermission(role.ViewOwnFiles), h.LookupMyImages)

	groupAdmin := router.Group("/image/admin")
	groupAdmin.Get("/list", middleware.RequirePermission(role.ViewOtherFiles), h.LookupAllImages)
	groupAdmin.Get("/list/:id", middleware.RequirePermission(role.ViewOtherFiles), h.LookupImagesByUserID)
}
