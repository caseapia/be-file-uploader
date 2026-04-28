package file

import (
	"strconv"

	"be-file-uploader/internal/models/requests"
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/file"
	"be-file-uploader/internal/service/user"
	"be-file-uploader/pkg/enums/role"
	"be-file-uploader/pkg/utils/account"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	imageService *image.Service
	userService  *user.Service
	repository   *mysql.Repository
}

func NewHandler(imageService *image.Service, userService *user.Service, repository *mysql.Repository) *Handler {
	return &Handler{imageService: imageService, userService: userService, repository: repository}
}

func (h *Handler) InitUpload(ctx fiber.Ctx) error {
	uploader := account.GetUserFromContext(ctx)

	var req requests.InitUpload
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	resp, err := h.imageService.InitMultipartUpload(ctx.Context(), uploader, req)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, resp)
}

func (h *Handler) UploadChunk(ctx fiber.Ctx) error {
	uploadID := ctx.FormValue("upload_id")
	key := ctx.FormValue("key")
	partNumberRaw := ctx.FormValue("part_number")

	partNumber, err := strconv.ParseInt(partNumberRaw, 10, 32)
	if err != nil || uploadID == "" || key == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ERR_INVALID_PARAMS")
	}

	fileHeader, err := ctx.FormFile("chunk")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ERR_CHUNK_MISSING")
	}

	eTag, err := h.imageService.UploadChunk(ctx.Context(), uploadID, key, int32(partNumber), fileHeader)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, &fiber.Map{"eTag": eTag})
}

func (h *Handler) CompleteUpload(ctx fiber.Ctx) error {
	uploader := account.GetUserFromContext(ctx)

	var req requests.CompleteUpload
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	fl, err := h.imageService.CompleteMultipartUpload(ctx.Context(), uploader, &req)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusCreated, fl)
}

func (h *Handler) DeleteImage(ctx fiber.Ctx) error {
	var req requests.DeleteImage
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	requester := account.GetUserFromContext(ctx)

	usedStorage, err := h.imageService.DeleteFile(ctx, req.ImageID, requester)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, &fiber.Map{
		"used_storage": usedStorage,
		"status":       "OK",
	})
}

func (h *Handler) LookupMyImages(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	images, err := h.repository.SearchOwnFiles(ctx, sender)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, images)
}

func (h *Handler) LookupAllImages(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	images, err := h.imageService.LookupAllFiles(ctx, sender)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, images)
}

func (h *Handler) LookupImagesByUserID(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)

	images, err := h.repository.SearchFilesByUserID(ctx, id)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, images)
}

func (h *Handler) AddInAlbum(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)
	var req requests.AddImageInAlbum
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	img, err := h.imageService.AddImageInAlbum(ctx, sender, req.ImageID, req.AlbumID)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, img)
}

func (h *Handler) RemoveFromAlbum(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)
	var req requests.RemoveImageFromAlbum
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	img, err := h.imageService.RemoveImageFromAlbum(ctx, sender, req.ImageID)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, img)
}

func (h *Handler) LikePost(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)

	state, err := h.imageService.ToggleLike(ctx, sender, id, true)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, state)
}

func (h *Handler) RemoveLikeFromPost(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)

	state, err := h.imageService.ToggleLike(ctx, sender, id, false)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, state)
}

func (h *Handler) DownloadImage(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)

	link, err := h.imageService.DownloadFile(ctx, sender, id)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusCreated, link)
}

func (h *Handler) AddComment(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	var req requests.AddCommentToPost
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	comment, err := h.imageService.AddComment(ctx, sender, req.PostID, req.Content)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusCreated, comment)
}

func (h *Handler) LookupPostByID(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)

	findedImage, err := h.imageService.FindFile(ctx, sender, id)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, findedImage)
}

func (h *Handler) ShareXUpload(ctx fiber.Ctx) error {
	token := ctx.FormValue("token")
	u, err := h.userService.AuthByToken(ctx.Context(), token)
	if err != nil {
		return err
	}

	if !u.HasPermission(role.FileUpload) {
		return fiber.NewError(fiber.StatusForbidden, "ERR_NO_PERMISSION")
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ERR_IMAGE_MISSING")
	}

	initReq := requests.InitUpload{
		OriginalName: fileHeader.Filename,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		Size:         fileHeader.Size,
		IsPrivate:    false,
	}

	initResp, err := h.imageService.InitMultipartUpload(ctx.Context(), u, initReq)
	if err != nil {
		return err
	}

	file, _ := fileHeader.Open()
	defer file.Close()

	eTag, err := h.imageService.UploadChunk(ctx.Context(), initResp.UploadID, initResp.Key, 1, fileHeader)
	if err != nil {
		return err
	}

	completeReq := requests.CompleteUpload{
		UploadID:     initResp.UploadID,
		Key:          initResp.Key,
		OriginalName: initReq.OriginalName,
		MimeType:     initReq.MimeType,
		Size:         initReq.Size,
		IsPrivate:    initReq.IsPrivate,
		Parts: []requests.Part{
			{PartNumber: 1, ETag: eTag},
		},
	}

	img, err := h.imageService.CompleteMultipartUpload(ctx.Context(), u, &completeReq)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"url": img.URL,
	})
}
