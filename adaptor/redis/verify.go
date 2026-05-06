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

	SetVerifyCode(ctx context.Context, mobile, sceneCode string, value interface{}, expire time.Duration) error
	GetVerifyCode(ctx context.Context, mobile, sceneCode string) (string, error)
	DelVerifyCode(ctx context.Context, mobile, sceneCode string) error

	SetAdminUserToken(ctx context.Context, userID int64, token string, tokenData string, expire time.Duration) error
	GetAdminUserToken(ctx context.Context, token string) (string, error)
	CleanToken(ctx context.Context, userId int64) error
	SetCustomerUserToken(ctx context.Context, userID int64, token string, tokenData string, expire time.Duration) error
	GetCustomerUserToken(ctx context.Context, token string) (string, error)
	CleanCustomerToken(ctx context.Context, userID int64) error

	IncrPasswordErr(ctx context.Context, mobile string, expire time.Duration) (int64, error)
	DeletePasswordErr(ctx context.Context, mobile string) error
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

func fmtVerifyVerifyCode(mobile, sceneCode string) string {
	return fmt.Sprintf("%s:verify:code:%s:%s", config.ServerName, mobile, sceneCode)
}

func fmtVerifyAdminUserToken(token string) string {
	return fmt.Sprintf("%s:admin:user:token:%s", config.ServerName, token)
}

func fmtVerifyPasswordErr(mobile string) string {
	return fmt.Sprintf("%s:admin:user:password:errorcount:%s", config.ServerName, mobile)
}

func fmtUserMapTokenAdminUser(userId int64) string {
	return fmt.Sprintf("%s:admin:token:user:%d", config.ServerName, userId)
}

func fmtVerifyCustomerUserToken(token string) string {
	return fmt.Sprintf("%s:customer:user:token:%s", config.ServerName, token)
}

func fmtUserMapTokenCustomerUser(userID int64) string {
	return fmt.Sprintf("%s:customer:token:user:%d", config.ServerName, userID)
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

func (v *Verify) SetVerifyCode(ctx context.Context, mobile, sceneCode string, value interface{}, expire time.Duration) error {
	redisKey := fmtVerifyVerifyCode(mobile, sceneCode)
	return v.redis.Set(redisKey, value, expire).Err()
}

func (v *Verify) GetVerifyCode(ctx context.Context, mobile, sceneCode string) (string, error) {
	redisKey := fmtVerifyVerifyCode(mobile, sceneCode)
	return v.redis.Get(redisKey).Result()
}

func (v *Verify) DelVerifyCode(ctx context.Context, mobile, sceneCode string) error {
	redisKey := fmtVerifyVerifyCode(mobile, sceneCode)
	return v.redis.Del(redisKey).Err()
}

func (v *Verify) SetAdminUserToken(_ context.Context, userID int64, token string, tokenData string, expire time.Duration) error {
	redisKey := fmtVerifyAdminUserToken(token)
	_, err := v.redis.Set(redisKey, tokenData, expire).Result()
	if err != nil {
		return err
	}
	userMapTokenKey := fmtUserMapTokenAdminUser(userID)
	return v.redis.Set(userMapTokenKey, token, expire).Err()
}

func (v *Verify) GetAdminUserToken(_ context.Context, token string) (string, error) {
	redisKey := fmtVerifyAdminUserToken(token)
	get, err := v.redis.Get(redisKey).Result()
	if err != nil {
		return "", err
	}
	return get, nil
}

func (v *Verify) CleanToken(ctx context.Context, userId int64) error {
	userMapTokenKey := fmtUserMapTokenAdminUser(userId)
	token, err := v.redis.Get(userMapTokenKey).Result()
	if err != nil {
		return err
	}
	redisKey := fmtVerifyAdminUserToken(token)
	return v.redis.Del(redisKey, userMapTokenKey).Err()
}

func (v *Verify) SetCustomerUserToken(_ context.Context, userID int64, token string, tokenData string, expire time.Duration) error {
	redisKey := fmtVerifyCustomerUserToken(token)
	if err := v.redis.Set(redisKey, tokenData, expire).Err(); err != nil {
		return err
	}
	userMapTokenKey := fmtUserMapTokenCustomerUser(userID)
	return v.redis.Set(userMapTokenKey, token, expire).Err()
}

func (v *Verify) GetCustomerUserToken(_ context.Context, token string) (string, error) {
	redisKey := fmtVerifyCustomerUserToken(token)
	return v.redis.Get(redisKey).Result()
}

func (v *Verify) CleanCustomerToken(_ context.Context, userID int64) error {
	userMapTokenKey := fmtUserMapTokenCustomerUser(userID)
	token, err := v.redis.Get(userMapTokenKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	redisKey := fmtVerifyCustomerUserToken(token)
	return v.redis.Del(redisKey, userMapTokenKey).Err()
}

func (v *Verify) IncrPasswordErr(_ context.Context, mobile string, expire time.Duration) (int64, error) {
	redisKey := fmtVerifyPasswordErr(mobile)
	pipe := v.redis.Pipeline()
	incr, err := pipe.Incr(redisKey).Result()
	if err != nil {
		return 0, err
	}
	if incr == 1 {
		pipe.Expire(redisKey, expire)
	}
	_, err = pipe.Exec()
	return incr, err
}

func (v *Verify) DeletePasswordErr(_ context.Context, mobile string) error {
	redisKey := fmtVerifyPasswordErr(mobile)
	return v.redis.Del(redisKey).Err()
}
