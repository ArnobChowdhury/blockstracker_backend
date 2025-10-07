package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterChangeRoutes(router *gin.Engine, changeHandler *handlers.ChangeHandler, authMiddleware *middleware.AuthMiddleware) {
	changeRoutes := router.Group("/changes")

	changeRoutes.Use(authMiddleware.Handle)
	changeRoutes.GET("/sync", changeHandler.SyncChanges)
}
