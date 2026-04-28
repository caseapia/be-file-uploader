package r2

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Storage struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewStorage(accessKey, secretKey, bucket, publicURL string) (*Storage, error) {
	endpoint := "https://storage.yandexcloud.net"

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
		config.WithRegion("ru-central1"),
	)
	if err != nil {
		return nil, fmt.Errorf("r2: failed to load config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &Storage{
		client:     client,
		bucketName: bucket,
		publicURL:  publicURL,
	}, nil
}

func (s *Storage) CreateMultipartUpload(ctx context.Context, key, mimetype string) (uploadID string, error error) {
	input := s3.CreateMultipartUploadInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(mimetype),
	}

	resp, err := s.client.CreateMultipartUpload(ctx, &input)
	if err != nil {
		return "", fmt.Errorf("r2: failed to create multipart upload: %w", err)
	}

	return *resp.UploadId, nil
}

func (s *Storage) UploadPart(ctx context.Context, key, uploadID string, partNumber int32, data []byte) (eTag string, err error) {
	input := s3.UploadPartInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(key),
		PartNumber:    aws.Int32(partNumber),
		UploadId:      aws.String(uploadID),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
	}

	resp, err := s.client.UploadPart(ctx, &input)
	if err != nil {
		return "", fmt.Errorf("r2: failed to upload part: %w", err)
	}

	return *resp.ETag, nil
}

func (s *Storage) CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []types.CompletedPart) (url string, err error) {
	input := s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(s.bucketName),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: parts,
		},
	}

	_, err = s.client.CompleteMultipartUpload(ctx, &input)
	if err != nil {
		return "", fmt.Errorf("r2: failed to complete multipart upload: %w", err)
	}

	return fmt.Sprintf("%s/%s", s.publicURL, key), nil
}

func (s *Storage) Upload(ctx context.Context, key, mimeType string, data []byte) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentType:   aws.String(mimeType),
		ContentLength: aws.Int64(int64(len(data))),
	})
	if err != nil {
		return "", fmt.Errorf("r2: upload failed: %w", err)
	}

	return fmt.Sprintf("%s/%s", s.publicURL, key), nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("r2: delete failed: %w", err)
	}

	return nil
}

func (s *Storage) ReadAll(reader io.Reader) ([]byte, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("r2: read failed: %w", err)
	}

	return data, nil
}
