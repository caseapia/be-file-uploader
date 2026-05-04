package file

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/internal/models/relations"
	"be-file-uploader/internal/models/requests"
	"be-file-uploader/pkg/enums/role"

	"github.com/gofiber/fiber/v3"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (s *Service) DeleteFile(ctx fiber.Ctx, imageID int, requester *models.User) (uploadLimit *int64, err error) {
	image, err := s.repo.SearchFileByID(ctx, imageID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "ERR_IMAGE_NOTFOUND")
	}

	if image.UploadedBy != requester.ID && !requester.HasPermission(role.ManageFiles) {
		return nil, fiber.NewError(fiber.StatusForbidden, "ERR_IMAGE_FORBIDDEN")
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		requester.UsedStorage -= image.Size
		if err = s.repo.DeleteFile(ctx.Context(), tx, image); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "ERR_IMAGE_DELETE")
		}

		if image.UploadedBy != requester.ID {
			s.notify.CreateNotification(ctx.Context(), image.UploadedBy, fmt.Sprintf("NOTIFY_POST_DELETED_BY_MODERATOR+%s", image.OriginalName))
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
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_DELETE_IMAGE")
	}

	return &requester.UsedStorage, nil
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

	if !image.CanAccess(sender) {
		return nil, fiber.NewError(fiber.StatusNotFound, "ERR_FILE_NOTFOUND")
	}

	image.ResolveURL(sender)

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
		s.notify.CreateNotification(ctx.Context(), image.UploadedBy, fmt.Sprintf("NOTIFY_IMAGE_LIKED+%s+%s", image.OriginalName, sender.Username))
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

		s.notify.CreateNotification(ctx.Context(), post.UploadedBy, fmt.Sprintf("NOTIFY_IMAGE_ADD_COMMENT+%s+%s", post.OriginalName, comment.Author.Username))
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

func (s *Service) GrantAccess(ctx fiber.Ctx, sender *models.User, fileID, target int) (state bool, err error) {
	post, err := s.FindFile(ctx, sender, fileID)
	if err != nil {
		return false, err
	}
	user, err := s.repo.LookupUserByID(ctx.Context(), target)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fiber.NewError(fiber.StatusNotFound, "ERR_USER_NOTFOUND")
		}
		return false, err
	}

	if sender.ID != post.UploadedBy {
		return false, fiber.NewError(fiber.StatusNotFound, "ERR_FORBIDDEN")
	}
	if post.IsPrivate == false {
		return false, fiber.NewError(fiber.StatusConflict, "ERR_NOT_AVAILABLE")
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		grant := models.FileGrants{
			UserID:      user.ID,
			GrantedByID: sender.ID,
			IsOwner:     false,
			FileID:      post.ID,
		}

		err = s.repo.GrantAccess(ctx.Context(), tx, grant)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Service) RemoveAccess(ctx fiber.Ctx, sender *models.User, fileID, target int) (state bool, err error) {
	post, err := s.FindFile(ctx, sender, fileID)
	if err != nil {
		return false, err
	}
	user, err := s.repo.LookupUserByID(ctx.Context(), target)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fiber.NewError(fiber.StatusNotFound, "ERR_USER_NOTFOUND")
		}
		return false, err
	}
	if post.UploadedBy != sender.ID {
		return false, fiber.NewError(fiber.StatusNotFound, "ERR_FORBIDDEN")
	}

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		err = s.repo.RemoveAccess(ctx.Context(), tx, user.ID, post.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Service) EditFileDetails(ctx fiber.Ctx, sender *models.User, req requests.EditFileDetails) (file *models.File, err error) {
	post, err := s.FindFile(ctx, sender, req.FileID)
	if err != nil {
		return nil, err
	}

	if post.UploadedBy != sender.ID {
		return nil, fiber.NewError(fiber.StatusNotFound, "ERR_FORBIDDEN")
	}

	post.OriginalName = req.FileName
	post.IsPrivate = req.IsPrivate

	err = s.repo.WithTx(ctx.Context(), func(tx bun.Tx) (err error) {
		_, err = s.repo.UpdateFile(ctx.Context(), tx, post)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return post, nil
}
