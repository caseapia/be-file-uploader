package file

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	group := router.Group("/storage")
	upload := group.Group("/upload")

	upload.Post("/init", middleware.RequirePermission(role.FileUpload), h.InitUpload)
	upload.Post("/chunk", middleware.RequirePermission(role.FileUpload), h.UploadChunk)
	upload.Post("/complete", middleware.RequirePermission(role.FileUpload), h.CompleteUpload)

	group.Post("/delete", middleware.RequirePermission(role.FileUpload), h.DeleteImage)
	group.Get("/list/:id", middleware.RequirePermission(role.ViewOtherFiles), h.LookupImagesByUserID)
	group.Get("/my", middleware.RequirePermission(role.ViewOwnFiles), h.LookupMyImages)
	group.Put("/album/put", middleware.RequirePermission(role.FileUpload), h.AddInAlbum)
	group.Delete("/album/delete", middleware.RequirePermission(role.FileUpload), h.RemoveFromAlbum)
	group.Get("/list", middleware.RequirePermission(role.ViewOtherFiles), h.LookupAllImages)
	group.Patch("/post/action/like/:id", middleware.RequirePermission(role.ViewOtherFiles), h.LikePost)
	group.Delete("/post/action/likeRemove/:id", middleware.RequirePermission(role.ViewOtherFiles), h.RemoveLikeFromPost)
	group.Get("/post/action/download/:id", middleware.RequirePermission(role.DownloadOthersFiles), h.DownloadImage)
	group.Post("/post/action/addComment", middleware.RequirePermission(role.ViewOtherFiles), h.AddComment)
	group.Get("/post/:id", middleware.RequirePermission(role.ViewOtherFiles), h.LookupPostByID)
}

func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	group := router.Group("/storage/upload")

	group.Post("/sharex", h.ShareXUpload)
}
