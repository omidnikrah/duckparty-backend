package routes

import (
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/gin-gonic/gin"
	_ "github.com/omidnikrah/duckparty-backend/docs"
	"github.com/omidnikrah/duckparty-backend/internal/config"
	"github.com/omidnikrah/duckparty-backend/internal/handler"
	"github.com/omidnikrah/duckparty-backend/internal/middleware"
	duckService "github.com/omidnikrah/duckparty-backend/internal/service/duck"
	userService "github.com/omidnikrah/duckparty-backend/internal/service/user"
	"github.com/omidnikrah/duckparty-backend/internal/storage"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, rdb *redis.Client, sesClient *ses.Client, s3Storage *storage.S3Storage, config *config.Config) {
	userSvc := userService.NewService(db, rdb, sesClient, config)
	duckSvc := duckService.NewService(db, userSvc, s3Storage)

	userHandler := handler.NewUserHandler(userSvc)
	duckHandler := handler.NewDuckHandler(duckSvc)

	apiRouter := router.Group(config.ApiPrefix)

	v1Router := apiRouter.Group("/v1")
	v1Router.Use(middleware.ValidationErrorMiddleware())

	v1Router.POST("/auth", middleware.RateLimit(middleware.AuthRateLimit), userHandler.Authenticate)
	v1Router.POST("/auth/verify", userHandler.AuthenticateVerify)

	authenticated := v1Router.Group("/")
	authenticated.Use(middleware.AuthMiddleware(config))

	authenticated.PUT("/user/change-name", userHandler.UpdateName)

	v1Router.GET("/ducks", duckHandler.GetDucksList)
	authenticated.POST("/duck", middleware.RateLimit(middleware.CreateRateLimit), duckHandler.CreateDuck)
	authenticated.PUT("/duck/:duckId/reaction/:reaction", duckHandler.ReactionToDuck)

	v1Router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
