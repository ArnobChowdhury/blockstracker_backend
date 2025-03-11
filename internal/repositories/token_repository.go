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

// invalidateTokensLua is a Lua script that atomically invalidates the access token,
// its associated refresh token, and cleans up the token mapping in Redis.
const invalidateTokensLua = `
  local refreshToken = redis.call("GET", KEYS[2])
  if refreshToken then
    redis.call("SET", KEYS[1], "invalid", ARGV[1])        -- Invalidate access token
    redis.call("SET", "refreshToken:" .. refreshToken, "invalid", ARGV[2]) -- Invalidate refresh token
    redis.call("DEL", KEYS[2])                             -- Remove access-to-refresh mapping
    return 1
  else
    return 0
  end
`

var invalidateTokensScript = redis.NewScript(invalidateTokensLua)

func (r *tokenRepository) InvalidateAccessAndRefreshTokens(accessToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := invalidateTokensScript.Run(ctx, r.client,
		[]string{
			"accessToken:" + accessToken,
			"accessToRefresh:" + accessToken,
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
