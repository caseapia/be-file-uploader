package file

import (
	"context"
	"mime/multipart"

	"be-file-uploader/internal/models"
	"be-file-uploader/internal/models/requests"
	"be-file-uploader/internal/service/upload"
	"be-file-uploader/pkg/utils/generate"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"
)

func (s *Service) validateUploadLimits(u *models.User, size int64) error {
	if size > upload.MaxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "ERR_IMAGE_TOO_LARGE")
	}
	if u.UsedStorage+size > u.UploadLimit {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "ERR_QUOTA_EXCEEDED")
	}
	return nil
}

func (s *Service) processImageFile(fh *multipart.FileHeader) ([]byte, string, string, error) {
	file, err := s.upload.DetectMultipartFile(fh, upload.MaxFileSize, upload.FileMimeExtensions)
	if err != nil {
		return nil, "", "", err
	}

	return file.Data, file.MimeType, file.Extension, nil
}

func (s *Service) generateStorageKey(userID int, imgID, ext string) string {
	return s.upload.GenerateKey("images", userID, imgID, ext, true)
}

func (s *Service) InitMultipartUpload(ctx context.Context, uploader *models.User, req requests.InitUpload) (*requests.InitUploadResponse, error) {
	if err := s.validateUploadLimits(uploader, req.Size); err != nil {
		return nil, err
	}

	ext, ok := upload.FileMimeExtensions[req.MimeType]
	if !ok {
		return nil, fiber.NewError(fiber.StatusBadRequest, "ERR_MIMETYPE")
	}

	imageID, _ := generate.FileID()
	r2Key := s.generateStorageKey(uploader.ID, imageID, ext)

	uploadID, err := s.upload.CreateMultipartUpload(ctx, r2Key, req.MimeType)
	if err != nil {
		return nil, err
	}

	return &requests.InitUploadResponse{
		UploadID: uploadID,
		Key:      r2Key,
	}, nil
}

func (s *Service) UploadChunk(ctx context.Context, uploadID, key string, partNumber int32, fh *multipart.FileHeader) (string, error) {
	data, err := s.upload.ReadMultipartFile(fh, upload.MaxFileSize)
	if err != nil {
		return "", err
	}

	eTag, err := s.upload.UploadPart(ctx, key, uploadID, partNumber, data)
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

	publicURL, err := s.upload.CompleteMultipartUpload(ctx, req.Key, req.UploadID, s3Parts)
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
