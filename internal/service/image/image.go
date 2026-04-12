package image

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/pkg/enums/role"
	"be-file-uploader/pkg/utils/generate"

	"github.com/gofiber/fiber/v3"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

var allowedMime = map[string]string{
	"image/jpeg": "image/jpeg",
	"image/png":  "image/png",
	"image/webp": "image/webp",
	"image/gif":  "image/gif",
}

const maxFileSize = 10 << 20

func (s *Service) UploadImage(ctx fiber.Ctx, uploader *models.User) (image *models.Image, err error) {
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		slog.WithData(slog.M{
			"fileHeader": fileHeader,
			"err":        err,
			"uploader":   uploader,
		}).Error("Error uploading image")
		return nil, fiber.NewError(fiber.StatusBadRequest, "ERR_IMAGE_MISSING")
	}

	if fileHeader.Size > maxFileSize {
		slog.WithData(slog.M{
			"fileHeader":  fileHeader.Size,
			"err":         err,
			"maxFileSize": maxFileSize,
		}).Error("Error uploading image (Image too large)")
		return nil, fiber.NewError(fiber.StatusRequestEntityTooLarge, "ERR_IMAGE_TOO_LARGE")
	}

	file, err := fileHeader.Open()
	if err != nil {
		slog.WithData(slog.M{
			"fileHeader": fileHeader,
			"err":        err,
		}).Error("Error open image")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_OPEN_IMAGE")
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			slog.WithData(slog.M{
				"fileHeader": fileHeader,
				"err":        err,
			}).Error("Error closing file")
		}
	}(file)

	data, err := s.storage.ReadAll(file)
	if err != nil {
		slog.WithData(slog.M{
			"fileHeader": fileHeader,
			"data":       data,
			"err":        err,
		}).Error("Error reading file")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_FILE_READING")
	}

	mimeType := http.DetectContentType(data)
	ext, ok := allowedMime[mimeType]
	if !ok {
		slog.WithData(slog.M{
			"fileHeader":  fileHeader,
			"data":        data,
			"mimeType":    mimeType,
			"allowedMime": allowedMime,
		})
		return nil, fiber.NewError(fiber.StatusUnsupportedMediaType, "ERR_IMAGE_UNSUPPORTED_TYPE")
	}

	imageID, err := generate.ImageID()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_IMAGE_ID_GEN")
	}

	// images/{userID}/{year-month}/{id}.ext
	r2Key := path.Join(
		"images",
		fmt.Sprintf("%d", uploader.ID),
		time.Now().Format("2006-01"),
		imageID+ext,
	)

	publicURL, err := s.storage.Upload(ctx.Context(), r2Key, mimeType, data)
	if err != nil {
		slog.WithData(slog.M{
			"fileHeader": fileHeader,
			"data":       data,
			"mimeType":   mimeType,
			"err":        err,
		})
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_IMAGE_UPLOAD")
	}

	image = &models.Image{
		R2Key:        r2Key,
		URL:          publicURL,
		OriginalName: fileHeader.Filename,
		MimeType:     mimeType,
		Size:         fileHeader.Size,
		UploadedBy:   uploader.ID,
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) error {
		if err = s.repo.CreateImage(ctx, tx, image); err != nil {
			slog.WithData(slog.M{
				"fileHeader": fileHeader,
				"data":       data,
				"mimeType":   mimeType,
				"err":        err,
			}).Error("Error creating image")
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (s *Service) DeleteImage(ctx fiber.Ctx, imageID int, requester *models.User) error {
	image, err := s.repo.SearchImageByID(ctx, imageID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "ERR_IMAGE_NOTFOUND")
	}

	if image.UploadedBy != requester.ID || !requester.HasPermission(role.ManageFiles) {
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
