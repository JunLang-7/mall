package token

import (
	"context"

	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/utils/logger"
	"go.uber.org/zap"
)

// GetLarkUserAccessToken 获取飞书用户 access token
func (s *Service) GetLarkUserAccessToken(ctx context.Context, appCode int32, code, redirectUrl, scope string, force bool) (*AccessToken, common.Errno) {
	token, err := s.getLarkUserAccessToken(ctx, appCode, code, redirectUrl, scope, force)
	if err != nil {
		logger.Error("GetLarkUserAccessToken get token failed", zap.Error(err), zap.Any("appCode", appCode), zap.Any("code", code))
		return nil, *common.ServerErr.WithErr(err)
	}
	return token, common.OK
}

// getLarkUserAccessToken 获取飞书用户 access token，支持强制刷新
func (s *Service) getLarkUserAccessToken(ctx context.Context, appCode int32, code, redirectUrl, scope string, force bool) (*AccessToken, error) {
	// 定义获取 token 的函数
	getTokenFunc := func() (*AccessToken, error) {
		token, err := s.lark.GetLarkUserAccessToken(ctx, appCode, code, redirectUrl, scope)
		if err != nil {
			logger.Error("getLarkUserAccessToken GetLarkUserAccessToken get token failed", zap.Error(err), zap.Any("appCode", appCode), zap.Any("code", code))
			return nil, err
		}
		return &AccessToken{
			Token:     token.AccessToken,
			ExpiresIn: token.ExpiresIn,
		}, nil
	}
	// Authorization code exchanges are per-user and short-lived, so they should not be shared via app-level cache.
	if code != "" {
		return getTokenFunc()
	}
	lockKey := s.lockTokenKeyFmt(appCode)
	cacheKey := s.cacheTokenKeyFmt(appCode)
	if !force {
		// 尝试从缓存获取 token，避免频繁调用飞书接口
		token, err := s.getToken(ctx, cacheKey)
		if err != nil {
			logger.Error("getLarkUserAccessToken get cache failed", zap.Error(err), zap.Any("appCode", appCode))
			return nil, err
		}
		if token != nil && token.Token != "" {
			return token, nil
		}
	}
	// 缓存不存在或强制刷新，获取新 token 并更新缓存
	token, err := s.updateToken(ctx, getTokenFunc, lockKey, cacheKey)
	if err != nil {
		logger.Error("getLarkUserAccessToken updateToken failed", zap.Error(err), zap.Any("appCode", appCode), zap.Any("code", code))
		return nil, err
	}
	if token == nil || token.Token == "" {
		return getTokenFunc()
	}
	return token, nil
}
