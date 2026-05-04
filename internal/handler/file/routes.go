package file

import (
	"be-file-uploader/internal/middleware"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) RegisterPrivateRoutes(router fiber.Router) {
	group := router.Group("/storage")
	upload := group.Group("/upload")
	action := group.Group("/action")

	upload.Post("/init", middleware.RequirePermission(role.FileUpload), h.InitUpload)
	upload.Post("/chunk", middleware.RequirePermission(role.FileUpload), h.UploadChunk)
	upload.Post("/complete", middleware.RequirePermission(role.FileUpload), h.CompleteUpload)

	action.Delete("/delete", middleware.RequirePermission(role.FileUpload), h.DeleteFile)
	group.Get("/list/:id", middleware.RequirePermission(role.ViewOtherFiles), h.LookupFilesByUserID)
	group.Get("/my", middleware.RequirePermission(role.ViewOwnFiles), h.LookupMyFiles)
	action.Put("/album/put", middleware.RequirePermission(role.FileUpload), h.AddInAlbum)
	action.Delete("/album/delete", middleware.RequirePermission(role.FileUpload), h.RemoveFromAlbum)
	group.Get("/list", middleware.RequirePermission(role.ViewOtherFiles), h.LookupAllFiles)
	action.Patch("/like/:id", middleware.RequirePermission(role.ViewOtherFiles), h.LikePost)
	action.Delete("/likeRemove/:id", middleware.RequirePermission(role.ViewOtherFiles), h.RemoveLikeFromPost)
	action.Get("/download/:id", middleware.RequirePermission(role.DownloadOthersFiles), h.DownloadFile)
	action.Post("/addComment", middleware.RequirePermission(role.ViewOtherFiles), h.AddComment)
	group.Get("/post/:id", middleware.RequirePermission(role.ViewOtherFiles), h.LookupPostByID)
	action.Put("/access/grant", h.GrantAccess)
	action.Delete("/access/remove", h.RemoveAccess)
	action.Patch("/update", h.EditFileDetails)
}

func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	group := router.Group("/storage/upload")

	group.Post("/sharex", h.ShareXUpload)
}
