package file

import (
	"context"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"

	"be-file-uploader/internal/models"
	"be-file-uploader/internal/models/requests"
	"be-file-uploader/pkg/utils/generate"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v3"
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
	if os.Getenv("APP_MODE") == "DEV" {
		return path.Join(
			"images/dev",
			strconv.FormatInt(int64(userID), 10),
			time.Now().Format("2006-01"),
			imgID+ext,
		)
	}

	return path.Join(
		"images",
		strconv.FormatInt(int64(userID), 10),
		time.Now().Format("2006-01"),
		imgID+ext,
	)
}

func (s *Service) InitMultipartUpload(ctx context.Context, uploader *models.User, req requests.InitUpload) (*requests.InitUploadResponse, error) {
	if err := s.validateUploadLimits(uploader, req.Size); err != nil {
		return nil, err
	}

	ext, ok := allowedMime[req.MimeType]
	if !ok {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_MIMETYPE")
	}

	imageID, _ := generate.FileID()
	r2Key := s.generateStorageKey(uploader.ID, imageID, ext)

	uploadID, err := s.storage.CreateMultipartUpload(ctx, r2Key, req.MimeType)
	if err != nil {
		return nil, err
	}

	return &requests.InitUploadResponse{
		UploadID: uploadID,
		Key:      r2Key,
	}, nil
}

func (s *Service) UploadChunk(ctx context.Context, uploadID, key string, partNumber int32, fh *multipart.FileHeader) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "ERR_OPEN_IMAGE")
	}
	defer file.Close()

	data, err := s.storage.ReadAll(file)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "ERR_FILE_READING")
	}

	eTag, err := s.storage.UploadPart(ctx, key, uploadID, partNumber, data)
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "ERR_UPLOAD_CHUNK")
	}

	return eTag, nil
}

func (s *Service) CompleteMultipartUpload(ctx context.Context, uploader *models.User, req *requests.CompleteUpload) (*models.File, error) {
	var s3Parts []types.CompletedPart
	for _, p := range req.Parts {
		s3Parts = append(s3Parts, types.CompletedPart{
			PartNumber: aws.Int32(p.PartNumber),
			ETag:       aws.String(p.ETag),
		})
	}

	publicURL, err := s.storage.CompleteMultipartUpload(ctx, req.Key, req.UploadID, s3Parts)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_COMPLETE_MULTIPART")
	}

	file := models.File{
		R2Key:        req.Key,
		URL:          publicURL,
		OriginalName: req.OriginalName,
		MimeType:     req.MimeType,
		Size:         req.Size,
		UploadedBy:   uploader.ID,
		IsPrivate:    req.IsPrivate,
	}

	err = s.repo.WithTx(ctx, func(tx bun.Tx) error {
		if err := s.repo.ReserveDiskSpace(ctx, tx, uploader, req.Size); err != nil {
			return err
		}
		file, err = s.repo.CreateFile(ctx, tx, &file)

		grant := models.FileGrants{
			FileID:      file.ID,
			IsOwner:     true,
			UserID:      uploader.ID,
			GrantedByID: uploader.ID,
		}

		err = s.repo.GrantAccess(ctx, tx, grant)

		return nil
	})
	if err != nil {
		_ = s.storage.Delete(ctx, req.Key)
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &file, nil
}
