// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"blockstracker_backend/config"
	"blockstracker_backend/handlers"
	"blockstracker_backend/internal/database"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/pkg/logger"
)

// Injectors from wire.go:

func InitializeAuthHandler() (*handlers.AuthHandler, error) {
	db := database.DBProvider()
	userRepository := repositories.NewUserRepository(db)
	sugaredLogger := logger.LoggerProvider()
	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		return nil, err
	}
	authHandler := handlers.NewAuthHandler(userRepository, sugaredLogger, authConfig)
	return authHandler, nil
}
