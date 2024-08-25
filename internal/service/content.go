package service

import (
	"autoshop/pkg/logging"
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"log/slog"
	"time"
)

type ContentService interface {
	UploadContent(content []byte) (*uuid.UUID, error)
}

type ContentStorage struct {
	minioClient *minio.Client
}

func NewContentStorage(minioClient *minio.Client) *ContentStorage {
	return &ContentStorage{minioClient: minioClient}
}

func (c *ContentStorage) UploadContent(content []byte) (*uuid.UUID, error) {
	const op = "ContentStorage.UploadContent"
	log := logging.CreateLoggerWithOp(op)

	contentId := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := c.minioClient.PutObject(ctx, "default", contentId.String(), bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Warn("Ошибка при попытке загрузки контента в минио", slog.String("err", err.Error()))

		return nil, err
	}

	return &contentId, nil
}
