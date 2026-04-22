package album

import (
	"strconv"

	"be-file-uploader/internal/models/requests"
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/album"
	"be-file-uploader/pkg/utils/account"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	albumService *album.Service
	repository   *mysql.Repository
}

func NewHandler(albumService *album.Service, repository *mysql.Repository) *Handler {
	return &Handler{albumService: albumService, repository: repository}
}

func (h *Handler) CreateAlbum(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	var req requests.CreateAlbum
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	createdAlbum, err := h.albumService.CreateAlbum(ctx, *sender, req.AlbumName, req.IsPrivate)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 201, createdAlbum)
}

func (h *Handler) LookupAlbum(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)

	alb, err := h.albumService.AlbumLookup(ctx, sender, id)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, alb)
}

func (h *Handler) DeleteAlbum(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	idStr := ctx.Params("id")
	id, _ := strconv.Atoi(idStr)

	state, err := h.albumService.DeleteAlbum(ctx, sender, id)
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, state)
}

func (h *Handler) AllAlbums(ctx fiber.Ctx) error {
	albums, err := h.repository.LookupAllAlbums(ctx.Context())
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, albums)
}
