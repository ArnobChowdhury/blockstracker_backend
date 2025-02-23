package main

import (
	"fmt"
	"net/http"

	_ "blockstracker_backend/docs"
	"blockstracker_backend/internal/validators"
	"blockstracker_backend/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// probably we can get rid of it
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Air hot reload is working!")
}

// @title Your API Title
// @version 1.0
// @description This is your API description.
// @host localhost:5000
// @BasePath /api/v1
func main() {
	validators.RegisterCustomValidators()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Example API route
	// r.GET("/ping", PingHandler)

	// Versioned API group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", PingHandler)
		routes.RegisterTaskRoutes(v1) // Register task routes under /api/v1
	}

	// routes.RegisterTaskRoutes(r)

	fmt.Println("Server started on :8080")
	r.Run(":8080")
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
