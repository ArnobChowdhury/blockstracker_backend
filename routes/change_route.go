package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterChangeRoutes(rg *gin.RouterGroup, changeHandler *handlers.ChangeHandler, authMiddleware *middleware.AuthMiddleware) {
	changeRoutes := rg.Group("/changes")

	changeRoutes.Use(authMiddleware.Handle)
	changeRoutes.Use(authMiddleware.RequirePremium)
	changeRoutes.GET("/sync", changeHandler.SyncChanges)
}
