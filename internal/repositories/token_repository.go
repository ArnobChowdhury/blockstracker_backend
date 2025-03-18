package repositories

import (
	"blockstracker_backend/internal/utils"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepository interface {
	InvalidateAccessAndRefreshTokens(accessToken string) (int64, error)
	StoreAccessTokenAndRefreshToken(accessToken, refreshToken string) error
	GetRefreshToken(accessToken string) (string, error)
}

type tokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) TokenRepository {
	return &tokenRepository{client: client}
}
func getKeyNames(accessToken string) (accessTokenKey, accessToRefreshKey string) {
	return "accessToken:" + accessToken, AccessToRefreshPrefix + accessToken
}

const AccessToRefreshPrefix = "accessToRefresh:"

func (r *tokenRepository) InvalidateAccessAndRefreshTokens(accessToken string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Del(ctx, AccessToRefreshPrefix+accessToken).Result()
}

func (r *tokenRepository) StoreAccessTokenAndRefreshToken(accessToken, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hashedRefreshToken := sha256.Sum256([]byte(refreshToken))
	hashedRefreshTokenString := hex.EncodeToString(hashedRefreshToken[:])

	return r.client.Set(ctx, AccessToRefreshPrefix+accessToken, hashedRefreshTokenString, utils.RefreshTokenExpiry).Err()
}

func (r *tokenRepository) GetRefreshToken(accessToken string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Get(ctx, AccessToRefreshPrefix+accessToken).Result()
}
