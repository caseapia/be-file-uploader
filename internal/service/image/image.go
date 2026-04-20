package image

import (
	"database/sql"
	"errors"
	"mime/multipart"
	"net/http"
	"path"
	"slices"
	"strconv"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/enums/role"
	"be-file-uploader/pkg/utils/generate"

	"github.com/gofiber/fiber/v3"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

var allowedMime = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

const maxFileSize = 50 << 20

func (s *Service) validateUploadLimits(u *models.User, size int64) error {
	if size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "ERR_IMAGE_TOO_LARGE")
	}
	if u.UsedStorage+size > u.UploadLimit {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "ERR_QUOTA_EXCEEDED")
	}
	return nil
}

func (s *Service) processImageFile(fh *multipart.FileHeader) ([]byte, string, string, error) {
	file, err := fh.Open()
	if err != nil {
		return nil, "", "", fiber.NewError(fiber.StatusInternalServerError, "ERR_OPEN_IMAGE")
	}
	defer file.Close()

	data, err := s.storage.ReadAll(file)
	if err != nil {
		return nil, "", "", fiber.NewError(fiber.StatusInternalServerError, "ERR_FILE_READING")
	}

	mimeType := http.DetectContentType(data)
	ext, ok := allowedMime[mimeType]
	if !ok {
		return nil, "", "", fiber.NewError(fiber.StatusUnsupportedMediaType, "ERR_IMAGE_UNSUPPORTED_TYPE")
	}

	return data, mimeType, ext, nil
}

func (s *Service) generateStorageKey(userID int, imgID, ext string) string {
	return path.Join(
		"images",
		strconv.FormatInt(int64(userID), 10),
		time.Now().Format("2006-01"),
		imgID+ext,
	)
}

func (s *Service) UploadImage(ctx fiber.Ctx, uploader *models.User, isPrivate bool) (*models.Image, error) {
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "ERR_IMAGE_MISSING")
	}

	if err := s.validateUploadLimits(uploader, fileHeader.Size); err != nil {
		return nil, err
	}

	data, mimeType, ext, err := s.processImageFile(fileHeader)
	if err != nil {
		return nil, err
	}

	imageID, _ := generate.ImageID()
	r2Key := s.generateStorageKey(uploader.ID, imageID, ext)

	publicURL, err := s.storage.Upload(ctx.Context(), r2Key, mimeType, data)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_IMAGE_UPLOAD")
	}

	image := &models.Image{
		R2Key: r2Key, URL: publicURL, OriginalName: fileHeader.Filename,
		MimeType: mimeType, Size: fileHeader.Size, UploadedBy: uploader.ID,
		IsPrivate: isPrivate,
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) error {
		if err := s.repo.ReserveDiskSpace(ctx.Context(), tx, uploader, fileHeader.Size); err != nil {
			return err
		}
		return s.repo.CreateImage(ctx, tx, image)
	})

	if err != nil {
		_ = s.storage.Delete(ctx.Context(), r2Key)
		return nil, err
	}

	return image, nil
}

func (s *Service) DeleteImage(ctx fiber.Ctx, imageID int, requester *models.User) error {
	image, err := s.repo.SearchImageByID(ctx, imageID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "ERR_IMAGE_NOTFOUND")
	}

	if image.UploadedBy != requester.ID && !requester.HasPermission(role.ManageFiles) {
		return fiber.NewError(fiber.StatusForbidden, "ERR_IMAGE_FORBIDDEN")
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		requester.UsedStorage -= image.Size
		if err = s.repo.DeleteImage(ctx.Context(), tx, image); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_IMAGE_DELETE")
		}

		if _, err = s.repo.UpdateUser(ctx.Context(), tx, requester, "used_storage"); err != nil {
			return err
		}

		if err = s.storage.Delete(ctx.Context(), image.R2Key); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "ERR_DELETE_IMAGE")
	}

	return nil
}

func (s *Service) LookupUserImages(ctx fiber.Ctx, user *models.User, requester *models.User) (images []models.Image, err error) {
	images, err = s.repo.SearchImagesByUserID(ctx, user.ID)

	if requester.ID != user.ID && !requester.HasPermission(role.ManageFiles) {
		images = slices.DeleteFunc(images, func(image models.Image) bool {
			return image.IsPrivate
		})
	}

	return images, err
}

func (s *Service) lookupImageAndAlbum(
	ctx fiber.Ctx,
	senderID, imageID int,
	albumID *int,
) (image *models.Image, album *models.Album, err error) {

	image, err = s.repo.SearchImageByID(ctx.Context(), imageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fiber.NewError(fiber.StatusNotFound, "ERR_IMAGE_NOTFOUND")
		}
		return nil, nil, err
	}

	if albumID != nil {
		album, err = s.repo.LookupAlbumByID(ctx.Context(), *albumID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil, fiber.NewError(fiber.StatusNotFound, "ERR_ALBUM_NOTFOUND")
			}
			return nil, nil, err
		}

		if senderID != album.CreatedByID {
			return nil, nil, fiber.NewError(fiber.StatusForbidden, "ERR_ALBUM_FORBIDDEN")
		}
	}

	if senderID != image.UploadedBy {
		return nil, nil, fiber.NewError(fiber.StatusForbidden, "ERR_IMAGE_FORBIDDEN")
	}

	return image, album, nil
}

func (s *Service) AddImageInAlbum(ctx fiber.Ctx, sender *models.User, imageID, albumID int) (image *models.Image, err error) {
	image, album, err := s.lookupImageAndAlbum(ctx, sender.ID, imageID, &albumID)
	if err != nil {
		return nil, err
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		image.AlbumID = &album.ID
		_, err = s.repo.UpdateImage(ctx.Context(), tx, image)
		if err != nil {
			slog.WithData(slog.M{
				"album": album,
				"user":  sender,
				"image": image,
				"err":   err,
			}).Error("repo.UpdateImage")
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_IMAGE_UPLOAD")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (s *Service) RemoveImageFromAlbum(ctx fiber.Ctx, sender *models.User, imageID int) (image *models.Image, err error) {
	image, _, err = s.lookupImageAndAlbum(ctx, sender.ID, imageID, nil)
	if err != nil {
		return nil, err
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		image.AlbumID = nil
		_, err = s.repo.UpdateImage(ctx.Context(), tx, image)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_IMAGE_UPLOAD")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (s *Service) LookupAllImages(ctx fiber.Ctx, sender *models.User) (images []models.Image, err error) {
	images, err = s.repo.SearchAllImages(ctx)
	if err != nil {
		return nil, err
	}

	if !sender.HasPermission(role.ManageFiles) {
		images = slices.DeleteFunc(images, func(img models.Image) bool {
			return img.IsPrivate
		})
	}

	return images, err
}
