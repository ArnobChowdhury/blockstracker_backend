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

func LoadAuthConfig() error {
	accessSecretKey, ok := os.LookupEnv("JWT_ACCESS_SECRET")
	if !ok {
		return fmt.Errorf(messages.ErrJWTAccessSecretNotFoundInEnvironment)
	}

	refreshSecretKey, ok := os.LookupEnv("JWT_REFRESH_SECRET")
	if !ok {
		return fmt.Errorf(messages.ErrJWTRefreshSecretNotFoundInEnvironment)
	}

	AuthSecrets = &AuthConfig{
		AccessSecret:  accessSecretKey,
		RefreshSecret: refreshSecretKey,
	}
	return nil
}

func AuthConfigProvider() *AuthConfig {
	return AuthSecrets
}
