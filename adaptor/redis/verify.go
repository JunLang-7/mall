package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/config"
	"github.com/go-redis/redis"
)

type IVerify interface {
	SetCaptchaKey(ctx context.Context, key string, val string, expire time.Duration) error
	GetCaptchaKey(ctx context.Context, key string) (string, error)
	SetCaptchaTicket(ctx context.Context, key string, val string, expire time.Duration) error
	GetCaptchaTicket(ctx context.Context, key string) (string, error)
}

type Verify struct {
	redis *redis.Client
}

func NewVerify(adaptor adaptor.IAdaptor) *Verify {
	return &Verify{
		redis: adaptor.GetRedis(),
	}
}

func fmtVerifyCaptchaKey(key string) string {
	return fmt.Sprintf("%s:captcha:%s", config.ServerName, key)
}

func fmtVerifyCaptchaTicket(key string) string {
	return fmt.Sprintf("%s:captcha:ticket:%s", config.ServerName, key)
}

func (v *Verify) SetCaptchaKey(_ context.Context, key string, val string, expire time.Duration) error {
	redisKey := fmtVerifyCaptchaKey(key)
	return v.redis.Set(redisKey, val, expire).Err()
}

func (v *Verify) GetCaptchaKey(_ context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaKey(key)
	get, err := v.redis.Get(redisKey).Result()
	if err != nil {
		return "", err
	}
	v.redis.Del(redisKey)
	return get, nil
}

func (v *Verify) SetCaptchaTicket(_ context.Context, key string, val string, expire time.Duration) error {
	redisKey := fmtVerifyCaptchaTicket(key)
	return v.redis.Set(redisKey, val, expire).Err()
}

func (v *Verify) GetCaptchaTicket(_ context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaTicket(key)
	get, err := v.redis.Get(redisKey).Result()
	if err != nil {
		return "", err
	}
	v.redis.Del(redisKey)
	return get, nil
}
