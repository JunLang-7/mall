package redis

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/go-redis/redis"
)

type IAccessToken interface {
	SetAccessToken(ctx context.Context, key string, val string, expiration time.Duration) error
	GetAccessToken(ctx context.Context, key string) (string, error)
}

type AccessToken struct {
	redis *redis.Client
}

func NewAccessToken(adaptor adaptor.IAdaptor) IAccessToken {
	return &AccessToken{redis: adaptor.GetRedis()}
}

// SetAccessToken 设置 access token
func (a *AccessToken) SetAccessToken(_ context.Context, key string, val string, expiration time.Duration) error {
	return a.redis.Set(key, val, expiration).Err()
}

// GetAccessToken 获取 access token
func (a *AccessToken) GetAccessToken(_ context.Context, key string) (string, error) {
	return a.redis.Get(key).Result()
}
