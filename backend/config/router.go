package config

import (
	"backend/middlewares"

	userRepository "backend/internal/user/repository"

	authServices "backend/internal/auth/services"
	userServices "backend/internal/user/services"

	authControllers "backend/internal/auth/controllers"
	userControllers "backend/internal/user/controllers"

	authRoutes "backend/internal/auth/routes"
	userRoutes "backend/internal/user/routes"

	_ "backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	// "github.com/swaggo/gin-swagger/swaggerFiles"
)

// https://medium.com/@i.4erkasov/go-echo-framework-ddd-cqrs-part-1-637595917b3b

// @title           MeetSpace API
// @version         1.0
// @description     MeetSpace API documentation.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
//
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT
func SetupRouter(config *AppConfig) *gin.Engine {
	router := gin.Default()

	// Middlewares
	router.Use(logger.SetLogger()) // Equivalent to logger middleware
	router.Use(requestid.New())    // Request ID middleware
	router.Use(cors.Default())     // CORS middleware
	router.Use(middlewares.LoggerMiddleware(config.Logger, config.JWT))

	// Repository
	userRepo := userRepository.NewUserRepository(config.DB, config.Logger)

	// Services
	userSrv := userServices.NewUserService(userRepo, config.Logger)
	authSrv := authServices.NewAuthService(config.Logger, config.JWT, userSrv)

	// Controllers
	authCtrl := authControllers.NewAuthController(config.Logger, authSrv)
	userCtrl := userControllers.NewUserController(config.Logger, userSrv)

	// Routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	authRoutes.AuthRoutes(router, authCtrl)
	userRoutes.UserRoutes(router, userCtrl)

	return router
}
