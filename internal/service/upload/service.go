package upload

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"

	r2 "be-file-uploader/pkg/storage"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v3"
)

const (
	MaxFileSize   = 4 * 1024 * 1024 * 1024
	MaxAvatarSize = 5 * 1024 * 1024
)

var FileMimeExtensions = map[string]string{
	"image/jpeg":                   ".jpg",
	"image/png":                    ".png",
	"image/webp":                   ".webp",
	"image/gif":                    ".gif",
	"application/pdf":              ".pdf",
	"text/plain":                   ".txt",
	"application/zip":              ".zip",
	"application/x-rar-compressed": ".rar",
	"application/x-7z-compressed":  ".7z",
}

var AvatarMimeExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type Service struct {
	storage *r2.Storage
}

type MultipartFile struct {
	Data      []byte
	MimeType  string
	Extension string
}

func NewService(storage *r2.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) DetectMultipartFile(fh *multipart.FileHeader, maxSize int64, allowedMime map[string]string) (*MultipartFile, error) {
	data, err := s.ReadMultipartFile(fh, maxSize)
	if err != nil {
		return nil, err
	}

	mtype := mimetype.Detect(data)
	mimeType := mtype.String()

	ext, ok := allowedMime[mimeType]
	if !ok {
		return nil, fiber.NewError(fiber.StatusBadRequest, "ERR_MIMETYPE")
	}

	return &MultipartFile{
		Data:      data,
		MimeType:  mimeType,
		Extension: ext,
	}, nil
}

func (s *Service) ReadMultipartFile(fh *multipart.FileHeader, maxSize int64) ([]byte, error) {
	if fh.Size > maxSize {
		return nil, fiber.NewError(fiber.StatusRequestEntityTooLarge, "ERR_IMAGE_TOO_LARGE")
	}

	file, err := fh.Open()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_OPEN_IMAGE")
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxSize+1))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "ERR_FILE_READING")
	}
	if int64(len(data)) > maxSize {
		return nil, fiber.NewError(fiber.StatusRequestEntityTooLarge, "ERR_IMAGE_TOO_LARGE")
	}

	return data, nil
}

func (s *Service) GenerateKey(root string, userID int, objectID, ext string, dated bool) string {
	parts := []string{root}
	if os.Getenv("APP_MODE") == "DEV" {
		parts = []string{root, "dev"}
	}

	parts = append(parts, strconv.FormatInt(int64(userID), 10))
	if dated {
		parts = append(parts, time.Now().Format("2006-01"))
	}
	parts = append(parts, objectID+ext)

	return path.Join(parts...)
}

func (s *Service) CreateMultipartUpload(ctx context.Context, key, mimeType string) (string, error) {
	return s.storage.CreateMultipartUpload(ctx, key, mimeType)
}

func (s *Service) UploadPart(ctx context.Context, key, uploadID string, partNumber int32, data []byte) (string, error) {
	return s.storage.UploadPart(ctx, key, uploadID, partNumber, data)
}

func (s *Service) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []types.CompletedPart) (string, error) {
	return s.storage.CompleteMultipartUpload(ctx, key, uploadID, parts)
}

func (s *Service) Upload(ctx context.Context, key, mimeType string, data []byte) (string, error) {
	return s.storage.Upload(ctx, key, mimeType, data)
}
