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
	StoreAccessTokenAndRefreshToken(accessToken, refreshToken string) error
}

type tokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) TokenRepository {
	return &tokenRepository{client: client}
}

const invalidateTokensLua = `
    local accessTokenPrefix = "accessToken:"
    local refreshTokenPrefix = "refreshToken:"
    local accessToRefreshPrefix = "accessToRefresh:"

    local refreshTokenKey = accessToRefreshPrefix .. KEYS[1]
    local refreshToken = redis.call("GET", refreshTokenKey)

    if refreshToken and refreshToken ~= false then
        redis.call("SET", accessTokenPrefix .. KEYS[1], "invalid", "EX", ARGV[1])
        redis.call("SET", refreshTokenPrefix .. refreshToken, "invalid", "EX", ARGV[2])
        redis.call("DEL", refreshTokenKey)
        return 1
    end

    return 0
`

var invalidateTokensScript = redis.NewScript(invalidateTokensLua)

func (r *tokenRepository) InvalidateAccessAndRefreshTokens(accessToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := invalidateTokensScript.Run(ctx, r.client,
		[]string{
			accessToken,
		},
		utils.AccessTokenExpiry.Seconds(),
		utils.RefreshTokenExpiry.Seconds(),
	).Int()
	if err != nil {
		return err
	}
	if result == 0 {
		return apperrors.ErrRedisKeyNotFound
	}
	return nil
}

func (r *tokenRepository) StoreAccessTokenAndRefreshToken(accessToken, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Set(ctx, "accessToRefresh:"+accessToken, refreshToken, utils.AccessTokenExpiry).Err()
}
