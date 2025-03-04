package config

import (
	"blockstracker_backend/messages"
	"fmt"
	"os"
)

type AuthConfig struct {
	AccessSecret  string
	RefreshSecret string
}

var AuthSecrets *AuthConfig

func LoadAuthConfig() (*AuthConfig, error) {
	accessSecretKey, ok := os.LookupEnv("JWT_ACCESS_SECRET")
	if !ok {
		return nil, fmt.Errorf(messages.ErrJWTAccessSecretNotFoundInEnvironment)
	}

	refreshSecretKey, ok := os.LookupEnv("JWT_REFRESH_SECRET")
	if !ok {
		return nil, fmt.Errorf(messages.ErrJWTRefreshSecretNotFoundInEnvironment)
	}

	return &AuthConfig{
		AccessSecret:  accessSecretKey,
		RefreshSecret: refreshSecretKey,
	}, nil
}
