package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/redis"
	"github.com/JunLang-7/mall/adaptor/rpc"
	"github.com/JunLang-7/mall/config"
	"github.com/JunLang-7/mall/utils/logger"
	redispkg "github.com/go-redis/redis"
	"github.com/gogf/gf/util/gconv"
	"go.uber.org/zap"
)

type AccessToken struct {
	Token     string `yaml:"token"`
	ExpiresIn int64  `yaml:"expires_in"`
}

type GetTokenFunc func() (*AccessToken, error)

type Service struct {
	conf        *config.Config
	locker      redis.ILocker
	accessToken redis.IAccessToken
	lark        rpc.ILark
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		conf:        adaptor.GetConf(),
		locker:      redis.NewLocker(adaptor),
		accessToken: redis.NewAccessToken(adaptor),
		lark:        rpc.NewLark(adaptor),
	}
}

// 存储access token的key
func (s *Service) cacheTokenKeyFmt(appCode int32) string {
	return fmt.Sprintf("%s:cachetoken:%d", config.ServerName, appCode)
}

// 分布式锁key
func (s *Service) lockTokenKeyFmt(appCode int32) string {
	return fmt.Sprintf("%s:locktoken:%d", config.ServerName, appCode)
}

// updateToken 获取新 token 并更新缓存，使用分布式锁避免缓存击穿
func (s *Service) updateToken(ctx context.Context, getToken GetTokenFunc, lockKey, cacheKey string) (*AccessToken, error) {
	// 获取分布式锁，确保只有一个请求能获取新 token 并更新缓存
	locked, err := s.locker.GetLock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	if locked {
		token, err := getToken()
		if err != nil {
			logger.Error("updateToken getToken failed", zap.Error(err))
			return nil, err
		}

		// 将新 token 存储到 Redis，设置过期时间
		err = s.accessToken.SetAccessToken(ctx, cacheKey, gconv.String(token), time.Duration(token.ExpiresIn)*time.Second)
		if err != nil {
			logger.Error("updateToken SetAccessToken failed", zap.Error(err))
			return nil, err
		}
		return token, nil
	}
	// 等待锁结束
	err = s.locker.AwaitLock(ctx, lockKey, time.Second*5)
	if err != nil {
		logger.Error("updateToken AwaitLock failed", zap.Error(err))
		return nil, err
	}
	logger.Debug("updateToken getCache")
	return s.getCache(ctx, cacheKey)
}

// getToken 从缓存获取 token，返回 nil 表示缓存不存在
func (s *Service) getToken(ctx context.Context, cacheKey string) (*AccessToken, error) {
	token, err := s.getCache(ctx, cacheKey)
	if err == nil {
		return token, nil
	}
	if !errors.Is(err, redispkg.Nil) {
		return nil, err
	}
	return nil, nil
}

// getCache 从 Redis 获取 token，返回 nil 表示缓存不存在
func (s *Service) getCache(ctx context.Context, cacheKey string) (*AccessToken, error) {
	cacheValue, err := s.accessToken.GetAccessToken(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	var token AccessToken
	if err = gconv.Struct(cacheValue, &token); err != nil {
		return nil, err
	}
	return &token, nil
}
