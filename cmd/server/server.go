package server

import (
	"context"
	"log/slog"

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

	sesClient, err := client.NewSESClient(config)
	if err != nil {
		panic("failed to initialize SES client: " + err.Error())
	}

	s3Storage, err := storage.NewS3Storage(config)
	if err != nil {
		panic("failed to initialize S3 storage: " + err.Error())
	}

	cronScheduler, err := client.NewCron(context.Background(), db, slog.Default())
	if err != nil {
		panic("failed to initialize cron: " + err.Error())
	}
	defer func() {
		if err := cronScheduler.Shutdown(); err != nil {
			slog.Default().Error("failed to shutdown cron scheduler", "error", err)
		}
	}()

	router := gin.Default()
	routes.SetupRoutes(router, db, rdb, sesClient, s3Storage, config)
	router.Run(":" + config.AppPort)
}
