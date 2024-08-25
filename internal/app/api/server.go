package api

import (
	"autoshop/internal/config"
	"autoshop/internal/controllers"
	"autoshop/internal/service"
	"autoshop/internal/storage"
	"autoshop/pkg/logging"
	"context"
	"database/sql"
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
	db, err := sql.Open("sqlite3", cfg.DbConfig.DbConnectionString)
	if err != nil {
		panic(err)
	}

	userStorage := storage.NewSqliteUserStorage(db)
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

	endpoint := "127.0.0.1:9000"
	accessKeyID := "BsdGKMYfWizG3KX5jXet"
	secretAccessKey := "TXVDFVbu4E10mfChKPFTTxprYd8JhB7Vbbdla2Im"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
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
