package repositories

import (
	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/utils"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepository interface {
	InvalidateAccessAndRefreshTokens(accessToken string) error
}

type tokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) TokenRepository {
	return &tokenRepository{client: client}
}

func (r *tokenRepository) InvalidateAccessAndRefreshTokens(accessToken string) error {
	script := redis.NewScript(`
		local refreshToken = redis.call("GET", KEYS[2])
		if refreshToken then
			redis.call("SET", KEYS[1], "invalid", ARGV[1]) -- Invalidate access token
			redis.call("SET", "refreshToken:" .. refreshToken, "invalid", ARGV[2]) -- Invalidate refresh token
			redis.call("DEL", KEYS[2]) -- Delete accessToRefresh:token entry
			return 1
		else
			return 0
		end
	`)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := script.Run(ctx, r.client, []string{
		"accessToken:" + accessToken,
		"accessToRefresh:" + accessToken,
	}, utils.AccessTokenExpiry.Seconds(), utils.RefreshTokenExpiry.Seconds()).Int()

	if err != nil {
		return err
	}
	if result == 0 {
		return apperrors.ErrRedisKeyNotFound
	}

	return nil
}
