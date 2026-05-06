package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/JunLang-7/mall/adaptor/rpc"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/JunLang-7/mall/utils/tools"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// GetSmsCodeVerify 获取短信验证码
func (s *Service) GetSmsCodeVerify(ctx context.Context, req *dto.GetSmsCodeVerifyReq) common.Errno {
	// 从 Redis 中获取验证码数据
	_, err := s.verify.GetCaptchaTicket(ctx, req.Ticket)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return common.InvalidCaptchaErr
		}
		logger.Error("MobilePasswordLogin GetCaptchaTicket error", zap.Error(err), zap.String("mobile", req.Mobile))
		return *common.RedisErr.WithErr(err)
	}
	// 生成短信验证码
	verifyCode := tools.GenValidateCode(4)
	// 发送飞书消息通知
	tokenFunc := func(ctx context.Context, force bool) (string, error) {
		token, errno := s.token.GetLarkTenantAccessToken(ctx, consts.LarkAppCode, force)
		if !errno.IsOK() {
			return "", errors.New(common.ServerErr.ErrMsg)
		}
		return token.Token, nil
	}
	err = s.lark.SendLarkMsg(ctx, tokenFunc, &do.SendLarkMsg{
		AppCode: consts.LarkAppCode,
		OpenID:  s.conf.BizConf.LarkGroupID,
		IDType:  rpc.LarkChatGroupType,
		Content: fmt.Sprintf("<b>手机验证码通知</b>\\n\\n手机号：%s\\n验证码：%s", req.Mobile, verifyCode),
	})
	if err != nil {
		logger.Error("GetSmsCodeVerify SendLarkMsg error", zap.Error(err), zap.String("mobile", req.Mobile))
		return *common.ServerErr.WithErr(err)
	}
	// 将验证码存储到 Redis，并设置过期时间
	scene := req.Scene
	switch req.Scene {
	case "reset_password":
		scene = consts.CustomerResetPasswordSmsCode
	case "change_password":
		scene = consts.CustomerChangePasswordSmsCode
	default:
		scene = consts.CustomerMobileLoginSmsCode
	}
	err = s.verify.SetVerifyCode(ctx, req.Mobile, scene, verifyCode, consts.ExpireVerifyCodeErrTime)
	if err != nil {
		logger.Error("GetSmsCodeVerify SetVerifyCode error", zap.Error(err), zap.String("mobile", req.Mobile))
		return *common.RedisErr.WithErr(err)
	}
	return common.OK
}

// checkSmsVerifyCode 校验短信验证码
func (s *Service) checkSmsVerifyCode(ctx context.Context, mobile, scene, verifyCode string) bool {
	getCode, err := s.verify.GetVerifyCode(ctx, mobile, scene)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false
		}
		logger.Error("checkSmsVerifyCode GetVerifyCode error", zap.Error(err), zap.String("mobile", mobile))
	}
	return getCode == verifyCode
}
