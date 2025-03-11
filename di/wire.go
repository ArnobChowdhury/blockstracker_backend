//go:build wireinject
// +build wireinject

package di

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"blockstracker_backend/config"
	"blockstracker_backend/internal/database"
	"blockstracker_backend/internal/redis"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/pkg/logger"

	"github.com/google/wire"
)

func InitializeAuthHandler() (*handlers.AuthHandler, error) {
	wire.Build(
		database.DBProvider,
		repositories.NewUserRepository,
		logger.LoggerProvider,
		config.LoadAuthConfig,
		config.LoadRedisConfig,
		redis.NewRedisClient,
		repositories.NewTokenRepository,
		handlers.NewAuthHandler)
	return &handlers.AuthHandler{}, nil

}
func InitializeAuthMiddleware() (*middleware.AuthMiddleware, error) {
	wire.Build(
		logger.LoggerProvider,
		config.LoadAuthConfig,
		middleware.NewAuthMiddleware,
	)
	return &middleware.AuthMiddleware{}, nil
}
