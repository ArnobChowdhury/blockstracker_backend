package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterBillingRoutes(rg *gin.RouterGroup, billingHandler *handlers.BillingHandler, authMiddleware *middleware.AuthMiddleware) {
	changeRoutes := rg.Group("/billing")

	changeRoutes.Use(authMiddleware.Handle)
	changeRoutes.POST("/google/verify", billingHandler.Verify)
}
