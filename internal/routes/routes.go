package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/omidnikrah/duckparty-backend/internal/handler"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	userHandler := handler.NewUserHandler(db)
	duckHandler := handler.NewDuckHandler(db, userHandler)

	router.POST("/duck", duckHandler.CreateDuck)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
}
