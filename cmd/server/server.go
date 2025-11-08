package server

import (
	"github.com/gin-gonic/gin"
	"github.com/omidnikrah/duckparty-backend/internal/client"
	"github.com/omidnikrah/duckparty-backend/internal/config"
	"github.com/omidnikrah/duckparty-backend/internal/database"
	"github.com/omidnikrah/duckparty-backend/internal/routes"
	"github.com/omidnikrah/duckparty-backend/internal/storage"
)

func Setup() {
	config, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	db, err := database.Init(config)
	if err != nil {
		panic("failed to init database: " + err.Error())
	}

	defer database.Close(db)

	rdb := client.NewRedisClient(config)

	defer rdb.Close()

	s3Storage, err := storage.NewS3Storage(config)
	if err != nil {
		panic("failed to initialize S3 storage: " + err.Error())
	}

	router := gin.Default()
	routes.SetupRoutes(router, db, rdb, s3Storage, config)
	router.Run(":" + config.AppPort)
}
