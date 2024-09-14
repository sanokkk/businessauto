package api

import (
	"autoshop/internal/config"
	"autoshop/internal/controllers"
	"autoshop/internal/service"
	"autoshop/internal/storage"
	"autoshop/pkg/logging"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"log/slog"
	"time"
)

type Server struct {
	Handler
}

type Handler interface {
	Start(apiConfig config.ApiConfig)
}

func NewServer(cfg *config.Config, logger *slog.Logger) *Server {
	userStorage := storage.NewUsersStorage(&cfg.DbConfig)
	productStorage := storage.NewProductStore(&cfg.DbConfig)

	authService := service.NewJwtAuthService(userStorage)
	productSercice := service.NewProductService(productStorage)
	contentService := service.NewContentStorage(configureMinioClient(cfg))

	handler := controllers.NewHttpHandler(authService, productSercice, contentService, logger)

	return &Server{handler}
}

func configureMinioClient(cfg *config.Config) *minio.Client {
	logger := logging.CreateLoggerWithOp("Server.configureMinioClient")
	logger.Info("Начинаю конфигурировать клиента для minio")

	if !cfg.ContentConfig.UseContentStorage {
		return nil
	}

	minioCfg := cfg.ContentConfig

	endpoint := fmt.Sprintf("%s:%s", minioCfg.Host, minioCfg.Port)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.User, minioCfg.Secret, ""),
		Secure: minioCfg.UseSsl,
	})
	if err != nil {
		logger.Warn(err.Error())

		panic(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now setup

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	exist, err := minioClient.BucketExists(ctx, "default")
	if err != nil {
		logger.Warn(err.Error())

		panic(err)
	}

	if !exist {
		log.Println("Дефолт бакета нет, создаю")
		if err := minioClient.MakeBucket(context.Background(), "default", minio.MakeBucketOptions{}); err != nil {
			logger.Warn(err.Error())

			panic(err)
		}
	}

	return minioClient
}
