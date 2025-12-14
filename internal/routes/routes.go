package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/omidnikrah/duckparty-backend/docs"
	"github.com/omidnikrah/duckparty-backend/internal/config"
	"github.com/omidnikrah/duckparty-backend/internal/handler"
	"github.com/omidnikrah/duckparty-backend/internal/middleware"
	duckService "github.com/omidnikrah/duckparty-backend/internal/service/duck"
	userService "github.com/omidnikrah/duckparty-backend/internal/service/user"
	"github.com/omidnikrah/duckparty-backend/internal/storage"
	ws "github.com/omidnikrah/duckparty-backend/internal/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/resend/resend-go/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, rdb *redis.Client, resendClient *resend.Client, s3Storage *storage.S3Storage, config *config.Config, broadcaster *ws.SocketBroadcaster) {
	userSvc := userService.NewService(db, rdb, resendClient, config)
	duckSvc := duckService.NewService(db, userSvc, s3Storage, broadcaster)

	userHandler := handler.NewUserHandler(userSvc)
	duckHandler := handler.NewDuckHandler(duckSvc)
	wsHandler := handler.NewWebSocketHandler(broadcaster)

	apiRouter := router.Group(config.ApiPrefix)

	router.GET("/ws", wsHandler.HandleWebSocket)

	v1Router := apiRouter.Group("/v1")
	v1Router.Use(middleware.ValidationErrorMiddleware())

	v1Router.POST("/auth", middleware.RateLimit(middleware.AuthRateLimit), userHandler.Authenticate)
	v1Router.POST("/auth/verify", middleware.RateLimit(middleware.AuthRateLimit), userHandler.AuthenticateVerify)
	v1Router.POST("/auth/anonymous", middleware.RateLimit(middleware.AuthRateLimit), userHandler.CreateAnonymousUser)

	authenticated := v1Router.Group("/")
	authenticated.Use(middleware.AuthMiddleware(config))

	authenticated.PUT("/user/change-name", userHandler.UpdateName)
	authenticated.POST("/user/set-email", middleware.RateLimit(middleware.AuthRateLimit), userHandler.SetEmail)
	authenticated.POST("/user/verify-set-email", middleware.RateLimit(middleware.AuthRateLimit), userHandler.VerifySetEmail)
	authenticated.GET("/user", userHandler.GetMeUser)

	v1Router.GET("/user/:userId/ducks", duckHandler.GetUserDucks)
	v1Router.GET("/leaderboard", duckHandler.GetDucksLeaderboard)
	v1Router.GET("/ducks", duckHandler.GetDucksList)
	authenticated.POST("/duck", middleware.RateLimit(middleware.CreateRateLimit), duckHandler.CreateDuck)
	authenticated.PUT("/duck/:duckId/reaction/:reaction", duckHandler.ReactionToDuck)
	authenticated.DELETE("/duck/:duckId", duckHandler.RemoveDuck)

	v1Router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
