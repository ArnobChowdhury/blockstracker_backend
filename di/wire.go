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

func InitializeTaskHandler() (*handlers.TaskHandler, error) {
	wire.Build(
		database.DBProvider,
		repositories.NewTaskRepository,
		repositories.NewChangeRepository,
		logger.LoggerProvider,
		handlers.NewTaskHandler,
	)
	return &handlers.TaskHandler{}, nil
}

func InitializeTagHandler() (*handlers.TagHandler, error) {
	wire.Build(
		database.DBProvider,
		repositories.NewTagRepository,
		repositories.NewChangeRepository,
		logger.LoggerProvider,
		handlers.NewTagHandler,
	)
	return &handlers.TagHandler{}, nil
}

func InitializeSpaceHandler() (*handlers.SpaceHandler, error) {
	wire.Build(
		database.DBProvider,
		repositories.NewSpaceRepository,
		repositories.NewChangeRepository,
		logger.LoggerProvider,
		handlers.NewSpaceHandler,
	)
	return &handlers.SpaceHandler{}, nil
}

func InitializeChangeHandler() (*handlers.ChangeHandler, error) {
	wire.Build(
		database.DBProvider,
		repositories.NewChangeRepository,
		repositories.NewTaskRepository,
		repositories.NewTagRepository,
		repositories.NewSpaceRepository,
		logger.LoggerProvider,
		handlers.NewChangeHandler,
	)
	return &handlers.ChangeHandler{}, nil
}

func InitializeBillingHandler() (*handlers.BillingHandler, error) {
	wire.Build(
		database.DBProvider,
		repositories.NewUserRepository,
		repositories.NewTokenRepository,
		config.LoadAuthConfig,
		config.LoadRedisConfig,
		redis.NewRedisClient,
		logger.LoggerProvider,
		handlers.NewBillingHandler,
	)
	return &handlers.BillingHandler{}, nil
}
