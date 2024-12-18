package controllers

import (
	//_ "autoshop/docs"
	"autoshop/internal/config"
	"autoshop/internal/middleware"
	"autoshop/internal/service"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

//	@title		HTTP api controller
//	@version	1.0

//	@host		localhost:8081
//	@BasePath	/api

// starts http handler
func (r *HttpHandler) Start(apiConfig config.ApiConfig) {
	const op = "HttpHandler.Start"
	log := r.logger.With(slog.String("op", op))

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	configureMiddleware(router)

	public := router.Group("/api")

	public.GET("/user", middleware.Authenticate(), r.GetMyUser)
	public.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	addUserRoutes(public, r)
	addProductRoutes(public, r)
	addContentRoutes(public, r)
	addCategoriesRoutes(public, r)

	log.Info("Запускаю HTTP сервер")

	addr := fmt.Sprintf("%s:%d", apiConfig.Host, apiConfig.Port)

	if err := router.Run(addr); err != nil {
		log.Error(err.Error())

		panic(err)
	}
}

func configureMiddleware(router *gin.Engine) {
	cfg := config.MustLoadConfig()

	if cfg.ApiConfig.EnableAnyOrigin {
		router.Use(cors.Default())
	}
}

func addUserRoutes(public *gin.RouterGroup, r *HttpHandler) {
	users := public.Group("/users")
	users.POST("/register", r.Register)
	users.POST("/login", r.Login)
	users.GET("/reauth", middleware.Authenticate(), r.Reauth)
	users.GET("/", middleware.Authenticate(), r.GetMyUser)
}
func addProductRoutes(public *gin.RouterGroup, r *HttpHandler) {
	products := public.Group("/products")
	products.POST("/get", r.GetProducts)
}

func addContentRoutes(public *gin.RouterGroup, r *HttpHandler) {
	products := public.Group("/content")

	products.Use(middleware.CheckFeatureFlag())
	products.POST("/", r.UploadFile)
	products.GET("/", r.DownloadFile)
}

func addCategoriesRoutes(public *gin.RouterGroup, r *HttpHandler) {
	categories := public.Group("/categories")

	categories.GET("/", middleware.Authenticate(), middleware.CheckForRole("admin"), r.GetCategories)
	categories.POST("/", middleware.Authenticate(), middleware.CheckForRole("admin"), r.HandleAddCategory)
}
