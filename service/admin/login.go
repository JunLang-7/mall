package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/rpc"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/JunLang-7/mall/utils/tools"
	"github.com/go-redis/redis"
	"github.com/gogf/gf/util/gconv"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	err = s.verify.SetVerifyCode(ctx, req.Mobile, consts.AdminUserMobileLoginSmsCode, verifyCode, consts.ExpireVerifyCodeErrTime)
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

// MobileVerifyLogin 手机号验证码登录
func (s *Service) MobileVerifyLogin(ctx context.Context, req *dto.MobileVerifyCodeLoginReq) (*dto.LoginResp, common.Errno) {
	pass := s.checkSmsVerifyCode(ctx, req.Mobile, consts.AdminUserMobileLoginSmsCode, req.VerifyCode)
	if !pass {
		return nil, common.InvalidSmsCodeErr
	}
	adminUser, err := s.adminUser.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		logger.Error("MobileVerifyLogin GetUserByMobile error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	if adminUser == nil || adminUser.Status != consts.IsEnable {
		return nil, common.AdminUserNotExistErr
	}
	return s.handleAdminLogin(ctx, adminUser)
}

// MobilePasswordLogin 手机号密码登录
func (s *Service) MobilePasswordLogin(ctx context.Context, req *dto.MobileLoginReq) (*dto.LoginResp, common.Errno) {
	// 从 Redis 中获取验证码数据
	_, err := s.verify.GetCaptchaTicket(ctx, req.Ticket)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, common.InvalidCaptchaErr
		}
		logger.Error("MobilePasswordLogin GetCaptchaTicket error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, *common.RedisErr.WithErr(err)
	}

	// 根据手机号获取用户信息
	adminUser, err := s.adminUser.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		logger.Error("MobilePasswordLogin GetUserByMobile error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	// 用户不存在
	if adminUser == nil || adminUser.Status != consts.IsEnable {
		return nil, common.InvalidPasswordErr
	}
	// 进行用户密码校验累计
	errCount, err := s.verify.IncrPasswordErr(ctx, req.Mobile, consts.ExpirePasswordErrTime)
	if err != nil {
		logger.Error("MobilePasswordLogin IncrPassword error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.InvalidPasswordErr
	}
	if errCount > consts.PasswordErrMaxCount {
		// 限制密码错误次数，比如10分钟内不能超过三次错误
		return nil, common.PasswordErrLimitErr
	}
	if adminUser.Password != req.Password {
		return nil, common.InvalidPasswordErr
	}
	// 登录成功，删除密码错误计数
	_ = s.verify.DeletePasswordErr(ctx, req.Mobile)

	return s.handleAdminLogin(ctx, adminUser)
}

// MobilePasswordReset 手机号重置密码
func (s *Service) MobilePasswordReset(ctx context.Context, req *dto.MobilePasswordResetReq) common.Errno {
	pass := s.checkSmsVerifyCode(ctx, req.Mobile, consts.AdminUserResetPasswordSmsCode, req.VerifyCode)
	if !pass {
		return common.InvalidSmsCodeErr
	}
	if req.Password != req.ConfirmPassword {
		return common.ConfirmPasswordErr
	}
	adminUser, err := s.adminUser.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.AdminUserNotExistErr
		}
		logger.Error("MobilePasswordReset GetUserByMobile error", zap.Error(err), zap.String("mobile", req.Mobile))
		return common.AdminUserNotExistErr
	}
	if adminUser == nil || adminUser.Status != consts.IsEnable {
		return common.AdminUserNotExistErr
	}
	err = s.adminUser.UpdateUserPassword(ctx, &do.UpdateUserPassword{
		ID:       adminUser.ID,
		Password: req.ConfirmPassword,
	})
	if err != nil {
		logger.Error("MobilePasswordReset UpdateUserPassword error", zap.Error(err), zap.String("mobile", req.Mobile))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

// LarkQrCodeLogin 飞书扫码登录
func (s *Service) LarkQrCodeLogin(ctx context.Context, req *dto.LarkQrCodeLoginReq) (*dto.LoginResp, common.Errno) {
	// 获取飞书用户 access token
	accessToken, errno := s.token.GetLarkUserAccessToken(ctx, req.AppCode, req.Code, req.RedirectUrl, "", false)
	if !errno.IsOK() {
		logger.Error("LarkQrCodeLogin GetLarkUserAccessToken error", zap.Any("req", req))
		return nil, common.ServerErr
	}
	// 通过 access token 获取飞书用户信息
	larkUserInfo, err := s.lark.GetLarkUserInfo(ctx, accessToken.Token)
	if err != nil {
		logger.Error("LarkQrCodeLogin GetLarkUserInfo error", zap.Error(err), zap.Any("req", req))
		return nil, *common.ServerErr.WithErr(err)
	}
	// 根据飞书用户信息中的 OpenID 获取对应的管理员用户信息
	adminUser, err := s.adminUser.GetUserByLarkOpenID(ctx, larkUserInfo.OpenID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("LarkQrCodeLogin GetUserByLarkOpenID error", zap.Error(err), zap.Any("req", req))
			return nil, common.InvalidLarkOpenIDErr
		}
		logger.Error("LarkQrCodeLogin GetUserByLarkOpenID error", zap.Error(err), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	if adminUser == nil || adminUser.Status != consts.IsEnable {
		return nil, common.AdminUserNotExistErr
	}

	return s.handleAdminLogin(ctx, adminUser)
}

// handleAdminLogin 处理管理员登录成功后的逻辑，包括生成 token 和构造响应数据
func (s *Service) handleAdminLogin(ctx context.Context, adminUser *model.AdminUser) (*dto.LoginResp, common.Errno) {
	adminUserDto := dto.AdminUserDto{
		UserID:     adminUser.ID,
		Name:       adminUser.Name,
		NickName:   adminUser.NickName,
		Sex:        adminUser.Sex,
		Status:     adminUser.Status,
		Mobile:     adminUser.Mobile,
		LarkOpenID: adminUser.LarkOpenID,
		UpdateAt:   adminUser.UpdateAt.UnixMilli(),
		CreateAt:   adminUser.CreateAt.UnixMilli(),
	}
	// NOTE: 可使用JWT
	tokenUuid := tools.UUIDHex()
	// 处理token
	err := s.processToken(ctx, tokenUuid, &adminUserDto)
	if err != nil {
		logger.Error("processToken error", zap.Error(err))
		return nil, *common.RedisErr.WithErr(err)
	}
	return &dto.LoginResp{
		Token: tokenUuid,
		User:  adminUserDto,
	}, common.OK
}

// processToken 处理登录成功后的 token 生成和存储
func (s *Service) processToken(ctx context.Context, token string, adminUser *dto.AdminUserDto) error {
	err := s.verify.SetAdminUserToken(ctx, adminUser.UserID, token, gconv.String(adminUser), consts.ExpireAdminUserTokenTime)
	if err != nil {
		logger.Error("SetAdminUserToken error", zap.Error(err), zap.String("mobile", adminUser.Mobile))
		return err
	}
	return nil
}
