//go:build wireinject
// +build wireinject

package di

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"blockstracker_backend/config"
	"blockstracker_backend/internal/database"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitializeAuthHandler() (*handlers.AuthHandler, error) {
	wire.Build(
		database.DBProvider,
		repositories.NewUserRepository,
		logger.LoggerProvider,
		config.LoadAuthConfig,
		handlers.NewAuthHandler)
	return &handlers.AuthHandler{}, nil

}
func InitializeAuthMiddleware() (gin.HandlerFunc, error) {
	wire.Build(
		logger.LoggerProvider,
		config.LoadAuthConfig,
		middleware.NewAuthMiddleware,
	)
	return nil, nil
}
