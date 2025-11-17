package main

import "github.com/omidnikrah/duckparty-backend/cmd/server"

// @title           Duck Party API
// @version         1.0
// @description     API for Duck Party backend service
// @termsOfService  http://swagger.io/terms/

// @host      localhost:4030
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	server.Setup()
}
