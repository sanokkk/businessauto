package controllers

import (
	"autoshop/internal/config"
	"autoshop/internal/middleware"
	"autoshop/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type HttpHandler struct {
	authService    service.AuthService
	productService service.ProductsService
	contentService service.ContentService
	logger         *slog.Logger
	validate       *validator.Validate
}

func NewHttpHandler(
	authService service.AuthService,
	productService service.ProductsService,
	contentService service.ContentService,
	logger *slog.Logger) *HttpHandler {
	return &HttpHandler{
		authService:    authService,
		productService: productService,
		contentService: contentService,
		logger:         logger,
		validate:       validator.New(validator.WithRequiredStructEnabled())}
}

func (r *HttpHandler) Start(apiConfig config.ApiConfig) {
	const op = "HttpHandler.Start"
	log := r.logger.With(slog.String("op", op))

	router := gin.Default()
	router.Use(gin.Recovery())

	public := router.Group("/api")

	public.GET("/user", middleware.Authenticate(), r.GetMyUser)

	addUserRoutes(public, r)
	addProductRoutes(public, r)
	addContentRoutes(public, r)

	log.Info("Запускаю HTTP сервер")

	addr := fmt.Sprintf("%s:%d", apiConfig.Host, apiConfig.Port)

	if err := router.Run(addr); err != nil {
		log.Error(err.Error())

		panic(err)
	}
}

func addUserRoutes(public *gin.RouterGroup, r *HttpHandler) {
	users := public.Group("/users")
	users.POST("/register", r.Register)
	users.POST("/login", r.Login)
	users.GET("/reauth", r.Reauth)
}
func addProductRoutes(public *gin.RouterGroup, r *HttpHandler) {
	products := public.Group("/products")
	products.POST("/get", r.GetProducts)
}

func addContentRoutes(public *gin.RouterGroup, r *HttpHandler) {
	products := public.Group("/content")

	products.Use(middleware.CheckFeatureFlag())
	products.POST("/post", r.UploadFile)
}
