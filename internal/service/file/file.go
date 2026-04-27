package image

import (
	"database/sql"
	"errors"
	"mime/multipart"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/internal/models/relations"
	"be-file-uploader/pkg/enums/role"
	"be-file-uploader/pkg/utils/generate"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v3"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

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

	mtype := mimetype.Detect(data)
	m := mtype.String()

	ext, ok := allowedMime[mtype.String()]
	if !ok {
		return nil, "", "", fiber.NewError(fiber.StatusInternalServerError, "ERR_MIMETYPE")
	}

	return data, m, ext, nil
}

func (s *Service) generateStorageKey(userID int, imgID, ext string) string {
	return path.Join(
		"images",
		strconv.FormatInt(int64(userID), 10),
		time.Now().Format("2006-01"),
		imgID+ext,
	)
}

func (s *Service) UploadFile(ctx fiber.Ctx, uploader *models.User, isPrivate bool) (*models.File, error) {
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

	image := &models.File{
		R2Key: r2Key, URL: publicURL, OriginalName: fileHeader.Filename,
		MimeType: mimeType, Size: fileHeader.Size, UploadedBy: uploader.ID,
		IsPrivate: isPrivate,
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) error {
		if err := s.repo.ReserveDiskSpace(ctx.Context(), tx, uploader, fileHeader.Size); err != nil {
			return err
		}
		return s.repo.CreateFile(ctx, tx, image)
	})

	if err != nil {
		_ = s.storage.Delete(ctx.Context(), r2Key)
		return nil, err
	}

	return image, nil
}

func (s *Service) DeleteFile(ctx fiber.Ctx, imageID int, requester *models.User) error {
	image, err := s.repo.SearchFileByID(ctx, imageID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "ERR_IMAGE_NOTFOUND")
	}

	if image.UploadedBy != requester.ID && !requester.HasPermission(role.ManageFiles) {
		return fiber.NewError(fiber.StatusForbidden, "ERR_IMAGE_FORBIDDEN")
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		requester.UsedStorage -= image.Size
		if err = s.repo.DeleteFile(ctx.Context(), tx, image); err != nil {
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

func (s *Service) LookupUserFiles(ctx fiber.Ctx, user *models.User, requester *models.User) (images []models.File, err error) {
	images, err = s.repo.SearchFilesByUserID(ctx, user.ID)

	if requester.ID != user.ID && !requester.HasPermission(role.ManageFiles) {
		images = slices.DeleteFunc(images, func(image models.File) bool {
			return image.IsPrivate
		})
	}

	for i := range images {
		images[i].ResolveURL(requester)
	}

	return images, err
}

func (s *Service) lookupFileAndAlbum(
	ctx fiber.Ctx,
	senderID, imageID int,
	albumID *int,
) (image *models.File, album *models.Album, err error) {

	image, err = s.repo.SearchFileByID(ctx.Context(), imageID)
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

	// _, err = s.repo.AddView(ctx, s.repo.DB, models.ImageViews{ImageID: image.ID, ViewerID: senderID})

	return image, album, nil
}

func (s *Service) AddImageInAlbum(ctx fiber.Ctx, sender *models.User, imageID, albumID int) (image *models.File, err error) {
	image, album, err := s.lookupFileAndAlbum(ctx, sender.ID, imageID, &albumID)
	if err != nil {
		return nil, err
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		image.AlbumID = &album.ID
		_, err = s.repo.UpdateFile(ctx.Context(), tx, image)
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

func (s *Service) RemoveImageFromAlbum(ctx fiber.Ctx, sender *models.User, imageID int) (image *models.File, err error) {
	image, _, err = s.lookupFileAndAlbum(ctx, sender.ID, imageID, nil)
	if err != nil {
		return nil, err
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		image.AlbumID = nil
		_, err = s.repo.UpdateFile(ctx.Context(), tx, image)
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

func (s *Service) LookupAllFiles(ctx fiber.Ctx, sender *models.User) (images []models.File, err error) {
	images, err = s.repo.SearchAllFiles(ctx)
	if err != nil {
		return nil, err
	}

	if !sender.HasPermission(role.ManageFiles) {
		images = slices.DeleteFunc(images, func(img models.File) bool {
			return img.IsPrivate
		})
	}

	for i := range images {
		images[i].ResolveURL(sender)
	}

	return images, err
}

func (s *Service) FindFile(ctx fiber.Ctx, sender *models.User, imageID int) (*models.File, error) {
	image, err := s.repo.SearchFileByID(ctx.Context(), imageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "ERR_IMAGE_NOTFOUND")
		}
		return nil, err
	}

	if (image.IsPrivate && sender.ID != image.UploadedBy) && !sender.HasPermission(role.ManageFiles) {
		return nil, fiber.NewError(fiber.StatusNotFound, "ERR_IMAGE_NOTFOUND")
	}

	if strings.Contains(image.MimeType, "image") && (sender.ID != image.UploadedBy || !sender.HasPermission(role.ManageFiles)) {
		image.URL = ""
	}

	return image, nil
}

func (s *Service) ToggleLike(ctx fiber.Ctx, sender *models.User, imageID int, add bool) (bool, error) {
	image, err := s.FindFile(ctx, sender, imageID)
	if err != nil {
		return false, err
	}

	like := models.FileLike{
		ImageID:  image.ID,
		AuthorID: sender.ID,
	}

	if add {
		return s.repo.AddLike(ctx.Context(), s.repo.DB, like)
	}
	return s.repo.RemoveLike(ctx.Context(), s.repo.DB, like)
}

func (s *Service) DownloadFile(ctx fiber.Ctx, sender *models.User, fileID int) (link string, err error) {
	file, err := s.FindFile(ctx, sender, fileID)
	if err != nil {
		return "", err
	}

	if file.IsPrivate == true && sender.ID != file.UploadedBy {
		return "", fiber.NewError(fiber.StatusNotFound, "ERR_IMAGE_NOTFOUND")
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		file.Downloads = file.Downloads + 1

		_, err = s.repo.UpdateFile(ctx.Context(), tx, file)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_IMAGE_UPLOAD")
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return file.URL, err
}

func (s *Service) AddComment(ctx fiber.Ctx, sender *models.User, image int, content string) (comment *models.FileComment, err error) {
	post, err := s.FindFile(ctx, sender, image)
	if err != nil {
		return nil, err
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		comment = &models.FileComment{
			AuthorID:  sender.ID,
			ImageID:   post.ID,
			CreatedAt: time.Now(),
			Content:   content,
			Author: relations.User{
				ID:       sender.ID,
				Username: sender.Username,
			},
		}

		comment, err = s.repo.AddComment(ctx.Context(), tx, comment)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}
