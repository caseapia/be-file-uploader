package image

import (
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/enums/role"
	"be-file-uploader/pkg/utils/generate"

	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"
)

var allowedMime = map[string]string{
	"image/jpeg": "image/jpeg",
	"image/png":  "image/png",
	"image/webp": "image/webp",
	"image/gif":  "image/gif",
}

const maxFileSize = 10 << 20

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

func (s *Service) UploadImage(ctx fiber.Ctx, uploader *models.User) (*models.Image, error) {
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
		if err = s.repo.DeleteImage(ctx.Context(), tx, image); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_IMAGE_DELETE")
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
