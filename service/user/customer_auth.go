package user

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/JunLang-7/mall/utils/secure"
	"github.com/JunLang-7/mall/utils/tools"
	goredis "github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *Service) MobileVerifyLogin(ctx context.Context, req *dto.MobileVerifyCodeLoginReq) (*dto.CustomerLoginResp, common.Errno) {
	if !s.checkSmsVerifyCode(ctx, req.Mobile, consts.CustomerMobileLoginSmsCode, req.VerifyCode) {
		return nil, common.InvalidSmsCodeErr
	}
	user, errno := s.findOrCreateMobileUser(ctx, req.Mobile)
	if !errno.IsOK() {
		return nil, errno
	}
	return s.handleCustomerLogin(ctx, user)
}

func (s *Service) MobilePasswordLogin(ctx context.Context, req *dto.MobileLoginReq) (*dto.CustomerLoginResp, common.Errno) {
	if _, err := s.verify.GetCaptchaTicket(ctx, req.Ticket); err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, common.InvalidCaptchaErr
		}
		return nil, *common.RedisErr.WithErr(err)
	}
	user, err := s.getUserByMobile(ctx, req.Mobile)
	if err != nil || user == nil || user.Status != consts.IsEnable {
		return nil, common.InvalidPasswordErr
	}
	if user.Password == "" || !secure.CheckPassword(req.Password, user.Password) {
		return nil, common.InvalidPasswordErr
	}
	return s.handleCustomerLogin(ctx, user)
}

func (s *Service) MobilePasswordReset(ctx context.Context, req *dto.MobilePasswordResetReq) common.Errno {
	if !s.checkSmsVerifyCode(ctx, req.Mobile, consts.CustomerResetPasswordSmsCode, req.VerifyCode) {
		return common.InvalidSmsCodeErr
	}
	if req.Password != req.ConfirmPassword {
		return common.ConfirmPasswordErr
	}
	user, err := s.getUserByMobile(ctx, req.Mobile)
	if err != nil || user == nil || user.Status != consts.IsEnable {
		return common.InvalidPasswordErr
	}
	hash, err := secure.HashPassword(req.Password)
	if err != nil {
		return *common.ServerErr.WithErr(err)
	}
	q := query.Use(s.db).User
	if _, err = q.WithContext(ctx).Where(q.ID.Eq(user.ID)).Update(q.Password, hash); err != nil {
		return *common.DataBaseErr.WithErr(err)
	}
	_ = s.verify.CleanCustomerToken(ctx, user.ID)
	return common.OK
}

func (s *Service) GetCustomerUserInfo(ctx context.Context, userID int64) (*dto.CustomerUserInfoDto, common.Errno) {
	info, err := s.buildCustomerUserInfo(ctx, userID)
	if err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return info, common.OK
}

func (s *Service) ChangePassword(ctx context.Context, userID int64, req *dto.ChangePasswordReq) (*dto.ChangePasswordResp, common.Errno) {
	if req.NewPassword != req.ConfirmPassword {
		return nil, common.ConfirmPasswordErr
	}
	q := query.Use(s.db).User
	user, err := q.WithContext(ctx).Where(q.ID.Eq(userID)).First()
	if err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	if user.Password != "" && !secure.CheckPassword(req.OldPassword, user.Password) {
		return nil, common.InvalidPasswordErr
	}
	hash, err := secure.HashPassword(req.NewPassword)
	if err != nil {
		return nil, *common.ServerErr.WithErr(err)
	}
	if _, err = q.WithContext(ctx).Where(q.ID.Eq(userID)).Update(q.Password, hash); err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	_ = s.verify.CleanCustomerToken(ctx, userID)
	return &dto.ChangePasswordResp{ReloginRequired: true}, common.OK
}

func (s *Service) GetCustomerUserByToken(ctx context.Context, token string) (*common.User, common.Errno) {
	data, err := s.verify.GetCustomerUserToken(ctx, token)
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, common.AuthErr
		}
		return nil, *common.RedisErr.WithErr(err)
	}
	user := &common.User{}
	if err = json.Unmarshal([]byte(data), user); err != nil {
		return nil, *common.ServerErr.WithErr(err)
	}
	q := query.Use(s.db).User
	dbUser, err := q.WithContext(ctx).Where(q.ID.Eq(user.UserID)).First()
	if err != nil || dbUser.Status != consts.IsEnable {
		return nil, common.AuthErr
	}
	return user, common.OK
}

func (s *Service) findOrCreateMobileUser(ctx context.Context, mobile string) (*model.User, common.Errno) {
	user, err := s.getUserByMobile(ctx, mobile)
	if err == nil && user != nil {
		if user.Status != consts.IsEnable {
			return nil, common.AuthErr
		}
		return user, common.OK
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	now := time.Now()
	mobileAES, err := secure.EncryptAESGCM(mobile, s.mobileAESKey())
	if err != nil {
		return nil, *common.ServerErr.WithErr(err)
	}
	mobileHash := secure.MobileSHA256(mobile, s.conf.Security.MobileSHA256Salt)
	user = &model.User{
		NickName:    "用户" + mobileTail(mobile),
		Sex:         0,
		Status:      consts.IsEnable,
		CreateAt:    now,
		UpdateAt:    now,
		LastLoginAt: now,
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return tx.Create(&model.MobileUser{
			UserID:       user.ID,
			MobileAes:    mobileAES,
			MobileSha256: mobileHash,
			CreateAt:     now,
			UpdateAt:     now,
		}).Error
	})
	if err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return user, common.OK
}

func (s *Service) getUserByMobile(ctx context.Context, mobile string) (*model.User, error) {
	mq := query.Use(s.db).MobileUser
	uq := query.Use(s.db).User
	mobileHash := secure.MobileSHA256(mobile, s.conf.Security.MobileSHA256Salt)
	mobileUser, err := mq.WithContext(ctx).Where(mq.MobileSha256.Eq(mobileHash)).First()
	if err != nil {
		return nil, err
	}
	return uq.WithContext(ctx).Where(uq.ID.Eq(mobileUser.UserID)).First()
}

func (s *Service) handleCustomerLogin(ctx context.Context, user *model.User) (*dto.CustomerLoginResp, common.Errno) {
	now := time.Now()
	uq := query.Use(s.db).User
	_, _ = uq.WithContext(ctx).Where(uq.ID.Eq(user.ID)).Update(uq.LastLoginAt, now)
	info, err := s.buildCustomerUserInfo(ctx, user.ID)
	if err != nil {
		return nil, *common.DataBaseErr.WithErr(err)
	}
	token := tools.UUIDHex()
	tokenUser := &common.User{UserID: user.ID, NickName: user.NickName, Status: user.Status}
	data, _ := json.Marshal(tokenUser)
	if err = s.verify.SetCustomerUserToken(ctx, user.ID, token, string(data), s.customerTokenTTL()); err != nil {
		logger.Error("SetCustomerUserToken error", zap.Error(err))
		return nil, *common.RedisErr.WithErr(err)
	}
	return &dto.CustomerLoginResp{Token: token, UserInfo: info}, common.OK
}

func (s *Service) buildCustomerUserInfo(ctx context.Context, userID int64) (*dto.CustomerUserInfoDto, error) {
	q := query.Use(s.db)
	user, err := q.User.WithContext(ctx).Where(q.User.ID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}
	info := &dto.CustomerUserInfoDto{
		User: &dto.CustomerUserDto{
			UserID:      user.ID,
			NickName:    user.NickName,
			CreateAt:    user.CreateAt.UnixMilli(),
			Sex:         user.Sex,
			Status:      user.Status,
			LastLoginAt: user.LastLoginAt.UnixMilli(),
			UpdateAt:    user.UpdateAt.UnixMilli(),
			HasPassword: user.Password != "",
		},
		AppUsers: make([]*dto.CustomerAppUserDto, 0),
	}
	if user.IconKey != "" {
		urlMap, _ := s.storage.GetPreviewUrl(ctx, &do.GetPreviewUrl{Keys: []string{user.IconKey}, Expire: 6})
		info.User.IconURL = urlMap[user.IconKey]
	}
	if mobileUser, err := q.MobileUser.WithContext(ctx).Where(q.MobileUser.UserID.Eq(userID)).First(); err == nil {
		mobile, _ := secure.DecryptAESGCM(mobileUser.MobileAes, s.mobileAESKey())
		info.MobileUser = &dto.CustomerMobileUserDto{Mobile: mobile, UserID: userID}
	}
	if wx, err := q.WechatUser.WithContext(ctx).Where(q.WechatUser.UserID.Eq(userID)).First(); err == nil {
		info.User.WechatBind = true
		info.WechatUser = &dto.CustomerWechatUserDto{UserID: userID, UnionID: wx.UnionID}
	}
	if apps, err := q.AppUser.WithContext(ctx).Where(q.AppUser.UserID.Eq(userID)).Find(); err == nil {
		for _, app := range apps {
			info.AppUsers = append(info.AppUsers, &dto.CustomerAppUserDto{OpenID: app.OpenID, UserID: app.UserID, AppCode: app.AppCode})
		}
	}
	return info, nil
}

func (s *Service) mobileAESKey() string {
	if s.conf.Security.MobileAESKey == "" {
		return "mall-default-mobile-aes-key"
	}
	return s.conf.Security.MobileAESKey
}

func mobileTail(mobile string) string {
	if len(mobile) <= 4 {
		return mobile
	}
	return mobile[len(mobile)-4:]
}

func (s *Service) customerTokenTTL() time.Duration {
	if s.conf.Security.CustomerTokenTTLSec > 0 {
		return time.Duration(s.conf.Security.CustomerTokenTTLSec) * time.Second
	}
	return consts.ExpireCustomerTokenTime
}
