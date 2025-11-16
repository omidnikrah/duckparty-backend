package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/omidnikrah/duckparty-backend/internal/config"
	"github.com/omidnikrah/duckparty-backend/internal/handler"
	"github.com/omidnikrah/duckparty-backend/internal/middleware"
	duckService "github.com/omidnikrah/duckparty-backend/internal/service/duck"
	userService "github.com/omidnikrah/duckparty-backend/internal/service/user"
	"github.com/omidnikrah/duckparty-backend/internal/storage"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, rdb *redis.Client, s3Storage *storage.S3Storage, config *config.Config) {
	userSvc := userService.NewService(db, rdb, config)
	duckSvc := duckService.NewService(db, userSvc, s3Storage)

	userHandler := handler.NewUserHandler(userSvc)
	duckHandler := handler.NewDuckHandler(duckSvc)

	v1Router := router.Group("/api")
	v1Router.Use(middleware.ValidationErrorMiddleware())

	v1Router.POST("/auth", middleware.RateLimit(middleware.AuthRateLimit), userHandler.Authenticate)
	v1Router.POST("/auth/verify", userHandler.AuthenticateVerify)

	authenticated := v1Router.Group("/")
	authenticated.Use(middleware.AuthMiddleware(config))

	v1Router.GET("/ducks", duckHandler.GetDucksList)
	authenticated.POST("/duck", middleware.RateLimit(middleware.CreateRateLimit), duckHandler.CreateDuck)
	authenticated.PUT("/duck/:duckId/reaction/:reaction", duckHandler.ReactionToDuck)

	v1Router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
}
