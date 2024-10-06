package controllers

import (
	//_ "autoshop/docs"
	"autoshop/internal/config"
	"autoshop/internal/middleware"
	"autoshop/internal/service"
	"fmt"
	cors "github.com/gofiber/fiber/v2/middleware/cors"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	//"github.com/gin-contrib/cors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
/*func (r *HttpHandler) Start(apiConfig config.ApiConfig) {
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
}*/

//	@title		HTTP api controller
//	@version	1.0

// @host		localhost:8080
// @BasePath	/api
func (r *HttpHandler) Start(apiConfig config.ApiConfig) {
	const op = "HttpHandler.Start"
	log := r.logger.With(slog.String("op", op))

	router := fiber.New()
	router.Use(cors.New())

	router.Use(recover2.New())

	public := router.Group("/api")

	public.Get("/swagger/*", swagger.HandlerDefault)

	addProductRoutes(public, r)
	addUserRoutes(public, r)
	addContentRoutes(public, r)
	addCategoriesRoutes(public, r)

	log.Info("Запускаю HTTP сервер")

	addr := fmt.Sprintf("%s:%d", apiConfig.Host, apiConfig.Port)

	if err := router.Listen(addr); err != nil {
		log.Error(err.Error())

		panic(err)
	}
}

/*func configureMiddleware(router *gin.Engine) {
	corsCfg := configureCors()
	router.Use(cors.New(corsCfg))
}*/

/*func configureCors() cors.Config {
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowCredentials = true
	corsCfg.AllowFiles = true
	corsCfg.AllowMethods = []string{"GET", "POST", "DELETE", "PATCH"}
	corsCfg.AllowAllOrigins = true
	return corsCfg
}*/

func addUserRoutes(public fiber.Router, r *HttpHandler) {
	users := public.Group("/users")
	users.Post("/register", r.Register)
	users.Post("/login", r.Login)
	users.Get("/reauth", middleware.AuthenticateFiber(), r.Reauth)
	users.Get("/", middleware.AuthenticateFiber(), r.GetMyUser)
}

func addProductRoutes(public fiber.Router, r *HttpHandler) {
	products := public.Group("/products")
	products.Post("/get", r.GetProducts)
}

func addContentRoutes(public fiber.Router, r *HttpHandler) {
	products := public.Group("/content")

	products.Use(middleware.CheckFeatureFlag())
	products.Post("/", r.UploadFile)
	products.Get("/", r.DownloadFile)
}

func addCategoriesRoutes(public fiber.Router, r *HttpHandler) {
	categories := public.Group("/categories")

	categories.Get("/", r.GetCategories)
	categories.Post("/", middleware.AuthenticateFiber(), middleware.CheckForRole("admin"), r.HandleAddCategory)
}
