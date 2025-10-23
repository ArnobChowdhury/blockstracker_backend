package config

import (
	"blockstracker_backend/messages"
	"fmt"
	"os"
)

type AuthConfig struct {
	AccessSecret          string
	RefreshSecret         string
	GoogleWebClientID     string
	GoogleWebClientSecret string
}

func LoadAuthConfig() (*AuthConfig, error) {
	accessSecretKey, ok := os.LookupEnv("JWT_ACCESS_SECRET")
	if !ok {
		return nil, fmt.Errorf(messages.ErrJWTAccessSecretNotFoundInEnvironment)
	}

	refreshSecretKey, ok := os.LookupEnv("JWT_REFRESH_SECRET")
	if !ok {
		return nil, fmt.Errorf(messages.ErrJWTRefreshSecretNotFoundInEnvironment)
	}

	googleWebClientId, ok := os.LookupEnv("GOOGLE_WEB_CLIENT_ID")
	if !ok {
		return nil, fmt.Errorf(messages.ErrGoogleWebClientIdNotFoundInEnvironment)
	}

	googleWebClientSecret, ok := os.LookupEnv("GOOGLE_WEB_CLIENT_SECRET")
	if !ok {
		return nil, fmt.Errorf(messages.ErrGoogleWebClientSecretNotFoundInEnvironment)
	}

	return &AuthConfig{
		AccessSecret:          accessSecretKey,
		RefreshSecret:         refreshSecretKey,
		GoogleWebClientID:     googleWebClientId,
		GoogleWebClientSecret: googleWebClientSecret,
	}, nil
}
