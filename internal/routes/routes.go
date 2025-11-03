package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/omidnikrah/duckparty-backend/internal/handler"
	duckService "github.com/omidnikrah/duckparty-backend/internal/service/duck"
	userService "github.com/omidnikrah/duckparty-backend/internal/service/user"
	"github.com/omidnikrah/duckparty-backend/internal/storage"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, s3Storage *storage.S3Storage) {
	userService := userService.NewService(db)
	duckService := duckService.NewService(db, userService, s3Storage)

	duckHandler := handler.NewDuckHandler(duckService)

	v1Router := router.Group("/api/")

	v1Router.POST("/duck", duckHandler.CreateDuck)
	v1Router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
}
