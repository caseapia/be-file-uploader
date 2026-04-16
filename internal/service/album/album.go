package album

import (
	"database/sql"
	"errors"
	"slices"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"
)

func (s *Service) CreateAlbum(ctx fiber.Ctx, creator models.User, albumName string, isPublic bool) (album *models.Album, err error) {
	album = &models.Album{
		CreatedByID: creator.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        albumName,
		Options: models.AlbumOptions{
			IsPublic: isPublic,
		},
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		album, err = s.repo.CreateAlbum(ctx.Context(), tx, *album)
		if err != nil {
			return err
		}

		return nil
	})

	return album, err
}

func (s *Service) AlbumLookup(ctx fiber.Ctx, sender *models.User, albumID int) (album *models.Album, err error) {
	album, err = s.repo.LookupAlbumByID(ctx.Context(), albumID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "ALBUM_NOTFOUND")
		}
		return nil, err
	}

	if !album.Options.IsPublic {
		if sender.ID != album.CreatedByID && !sender.HasPermission(role.ManageFiles) {
			return nil, fiber.NewError(fiber.StatusNotFound, "ALBUM_NOTFOUND")
		}
	}

	if sender.ID != album.CreatedByID && !sender.HasPermission(role.ManageFiles) {
		album.Items = slices.DeleteFunc(album.Items, func(img models.Image) bool {
			return img.IsPrivate
		})
	}

	return album, nil
}
