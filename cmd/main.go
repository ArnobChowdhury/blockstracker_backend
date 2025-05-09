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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @title Blockstracker
// @version 1.0
// @description Blockstracker API
// @host localhost:5000
// @BasePath /api/v1
func main() {
	defer logger.Log.Sync()

	validators.RegisterCustomValidators()
	database.ConnectDatabase()

	authHandler, err := di.InitializeAuthHandler()
	if err != nil {
		log.Fatalf("Error initializing auth handler: %s", err.Error())
	}
	authMiddleware, err := di.InitializeAuthMiddleware()
	if err != nil {
		log.Fatalf("Error initializing auth middleware: %s", err.Error())
	}

	taskHandler, err := di.InitializeTaskHandler()
	if err != nil {
		log.Fatalf("Error initializing task handler: %s", err.Error())
	}

	tagHandler, err := di.InitializeTagHandler()
	if err != nil {
		log.Fatalf("Error initializing tag handler: %s", err.Error())
	}

	spaceHandler, err := di.InitializeSpaceHandler()
	if err != nil {
		log.Fatalf("Error initializing space handler: %s", err.Error())
	}

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", PingHandler)

		routes.RegisterAuthRoutes(v1, authHandler, authMiddleware)
		routes.RegisterTaskRoutes(v1, taskHandler, authMiddleware)
		routes.RegisterTagRoutes(v1, tagHandler, authMiddleware)
		routes.RegisterSpaceRoutes(v1, spaceHandler, authMiddleware)
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
