package server

import (
	"context"
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/omidnikrah/duckparty-backend/internal/client"
	"github.com/omidnikrah/duckparty-backend/internal/config"
	"github.com/omidnikrah/duckparty-backend/internal/database"
	"github.com/omidnikrah/duckparty-backend/internal/routes"
	"github.com/omidnikrah/duckparty-backend/internal/storage"
	ws "github.com/omidnikrah/duckparty-backend/internal/websocket"
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

	resendClient := client.NewResendClient(config)

	r2Storage, err := storage.NewR2Storage(config)
	if err != nil {
		panic("failed to initialize R2 storage: " + err.Error())
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

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	broadcaster := ws.NewSocketBroadcaster()

	routes.SetupRoutes(router, db, rdb, resendClient, r2Storage, config, broadcaster)

	router.Run(":" + config.AppPort)
}
