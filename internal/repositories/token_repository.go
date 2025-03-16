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

const invalidateTokensLua = `
    local refreshTokenPrefix = "refreshToken:"
    local refreshToken = redis.call("GET", KEYS[2])

    if refreshToken and refreshToken ~= false then
        redis.call("SET", KEYS[1], "invalid", "EX", ARGV[1])
        redis.call("DEL", KEYS[2])
        return 1
    end

    return 0
`

var invalidateTokensScript = redis.NewScript(invalidateTokensLua)

func (r *tokenRepository) InvalidateAccessAndRefreshTokens(accessToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accessTokenKey, accessToRefreshKey := getKeyNames(accessToken)
	result, err := invalidateTokensScript.Run(ctx, r.client,
		[]string{
			accessTokenKey,
			accessToRefreshKey,
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

	return r.client.Set(ctx, AccessToRefreshPrefix+accessToken, refreshToken, utils.RefreshTokenExpiry).Err()
}

func (r *tokenRepository) GetRefreshToken(accessToken string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Get(ctx, AccessToRefreshPrefix+accessToken).Result()
}
