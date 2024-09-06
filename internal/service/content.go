package service

import (
	"autoshop/pkg/custom_errors"
	"autoshop/pkg/logging"
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"io"
	"log/slog"
	"time"
)

type ContentService interface {
	UploadContent(input FileInput) (*uuid.UUID, error)
	DownloadContent(contentId string) ([]byte, error)
}

type FileInput struct {
	Size    int64
	Content []byte
}

type ContentStorage struct {
	minioClient *minio.Client
}

func NewContentStorage(minioClient *minio.Client) *ContentStorage {
	return &ContentStorage{minioClient: minioClient}
}

func (c *ContentStorage) UploadContent(input FileInput) (*uuid.UUID, error) {
	const op = "ContentStorage.UploadContent"
	log := logging.CreateLoggerWithOp(op)

	contentId := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := c.minioClient.PutObject(ctx, "default", contentId.String(), bytes.NewReader(input.Content), input.Size, minio.PutObjectOptions{})
	if err != nil {
		log.Warn("Ошибка при попытке загрузки контента в минио", slog.String("err", err.Error()))

		return nil, err
	}

	return &contentId, nil
}

func (c *ContentStorage) DownloadContent(contentId string) ([]byte, error) {
	const op = "ContentStorage.DownloadContent"
	log := logging.CreateLoggerWithOp(op).With(slog.String("contentId", contentId))

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	content, err := c.minioClient.GetObject(ctx, "default", contentId, minio.GetObjectOptions{})
	stat, _ := content.Stat()
	if stat.Err != nil {
		log.Warn("Ошибка при получении информации о файле", slog.String("err", err.Error()))

		return nil, err
	}

	if stat.Size == 0 {
		log.Warn("Данный файл не существует")

		return nil, custom_errors.NoFileError
	}
	if err != nil {
		log.Warn("Ошибка при попытке скачивания контента", slog.String("err", err.Error()))

		return nil, err
	}

	fileContent, err := io.ReadAll(content)
	if err != nil {
		log.Warn("Ошибка при обработке контента", slog.String("err", err.Error()))

		return nil, err
	}

	return fileContent, nil
}
