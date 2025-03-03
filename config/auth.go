package config

import (
	"blockstracker_backend/messages"
	"fmt"
	"os"
)

type AuthConfig struct {
	AccessSecret string
	RefreshScret string
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
		AccessSecret: accessSecretKey,
		RefreshScret: refreshSecretKey,
	}
	return nil
}
