package server

import (
	"github.com/gin-gonic/gin"
	"github.com/omidnikrah/duckparty-backend/internal/config"
	"github.com/omidnikrah/duckparty-backend/internal/database"
	"github.com/omidnikrah/duckparty-backend/internal/routes"
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

	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run(":" + config.AppPort)
}