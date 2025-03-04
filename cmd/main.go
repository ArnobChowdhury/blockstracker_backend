package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"blockstracker_backend/di"
	_ "blockstracker_backend/docs"
	"blockstracker_backend/internal/database"
	"blockstracker_backend/internal/validators"
	"blockstracker_backend/pkg/logger"
	"blockstracker_backend/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Your API Title
// @version 1.0
// @description This is your API description.
// @host localhost:5000
// @BasePath /api/v1
func main() {
	defer logger.Log.Sync()

	// err := config.LoadAuthConfig()
	// if err != nil {
	// 	log.Fatalf("Error loading auth config: %v", err)
	// }

	validators.RegisterCustomValidators()
	database.ConnectDatabase()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authHandler, err := di.InitializeAuthHandler()
	if err != nil {
		log.Fatalf("Error loading auth config: %v", err)
	}
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", PingHandler)
		routes.RegisterAuthRoutes(v1, authHandler)
		routes.RegisterTaskRoutes(v1)
	}

	fmt.Println(strings.Repeat("🚀", 25))
	r.Run(":" + os.Getenv("PORT"))
}

// PingHandler example
// @Summary Ping example
// @Description Returns pong
// @Tags example
// @Success 200 {string} string "pong"
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
