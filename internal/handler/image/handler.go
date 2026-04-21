package image

import (
	"strconv"

	"be-file-uploader/internal/models/requests"
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/image"
	"be-file-uploader/pkg/utils/account"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	imageService *image.Service
	repository   *mysql.Repository
}

func NewHandler(imageService *image.Service, repository *mysql.Repository) *Handler {
	return &Handler{imageService: imageService, repository: repository}
}

func (h *Handler) UploadImage(ctx fiber.Ctx) error {
	uploader := account.GetUserFromContext(ctx)

	isPrivateRaw := ctx.FormValue("is_private")
	isPrivate, _ := strconv.ParseBool(isPrivateRaw)

	img, err := h.imageService.UploadImage(ctx, uploader, isPrivate)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 201, img)
}

func (h *Handler) DeleteImage(ctx fiber.Ctx) error {
	var req requests.DeleteImage
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	requester := account.GetUserFromContext(ctx)

	if err := h.imageService.DeleteImage(ctx, req.ImageID, requester); err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, "OK")
}

func (h *Handler) LookupMyImages(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	images, err := h.repository.SearchOwnImages(ctx, sender)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, images)
}

func (h *Handler) LookupAllImages(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	images, err := h.imageService.LookupAllImages(ctx, sender)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, images)
}

func (h *Handler) LookupImagesByUserID(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)

	images, err := h.repository.SearchImagesByUserID(ctx, id)
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

func (h *Handler) AddView(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)

	img, err := h.imageService.AddView(ctx, sender, id)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, img)
}
