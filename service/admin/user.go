package admin

import (
	"context"
	"errors"

	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/go-redis/redis"
	"github.com/gogf/gf/util/gconv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *Service) GetUserInfo(ctx context.Context, adminUser *common.AdminUser, userID int64) (*dto.AdminUserDto, common.Errno) {
	user, err := s.adminUser.GetUserInfo(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.InvalidPasswordErr
		}
		logger.Error("GetUserInfo error", zap.Error(err), zap.Int64("request_user_id", userID), zap.Any("admin_user", adminUser))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return &dto.AdminUserDto{UserID: user.ID, Name: user.Name}, common.OK
}

func (s *Service) CreateUser(ctx context.Context, adminUser *common.AdminUser, req *dto.CreateUserReq) (int64, common.Errno) {
	userId, err := s.adminUser.CreateUser(ctx, &do.CreateUser{
		AdminUserID: adminUser.UserID,
		Name:        req.Name,
		NickName:    req.NickName,
		Mobile:      req.Mobile,
		Sex:         req.Sex,
	})
	if err != nil {
		logger.Error("CreateUser error", zap.Error(err), zap.Any("req", req))
		return 0, *common.DataBaseErr.WithErr(err)
	}
	return userId, common.OK
}

func (s *Service) UpdateUser(ctx context.Context, adminUser *common.AdminUser, req *dto.UpdateUserReq) common.Errno {
	err := s.adminUser.UpdateUser(ctx, &do.UpdateUser{
		AdminUserID: adminUser.UserID,
		ID:          req.ID,
		Name:        req.Name,
		NickName:    req.NickName,
		Sex:         req.Sex,
	})
	if err != nil {
		logger.Error("UpdateUser error", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) UpdateUserStatus(ctx context.Context, adminUser *common.AdminUser, req *dto.UpdateUserStatusReq) common.Errno {
	err := s.adminUser.UpdateUserStatus(ctx, &do.UpdateUserStatus{
		AdminUserID: adminUser.UserID,
		ID:          req.ID,
		Status:      req.Status,
	})
	if err != nil {
		logger.Error("UpdateUserStatus error", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) DeleteUser(ctx context.Context, adminUser *common.AdminUser, req *dto.DeleteUserReq) common.Errno {
	err := s.adminUser.DeleteUser(ctx, req.ID)
	if err != nil {
		logger.Error("DeleteUser error", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) GetAdminUserByToken(ctx context.Context, token string) (*common.AdminUser, common.Errno) {
	userString, err := s.verify.GetAdminUserToken(ctx, token)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, common.AuthErr
		}
		logger.Error("GetAdminUserByToken redis get error", zap.Error(err), zap.Any("token", token))
		return nil, *common.RedisErr.WithErr(err)
	}
	adminUser := &common.AdminUser{}
	err = gconv.Struct(userString, adminUser)
	if err != nil {
		logger.Error("GetAdminUserByToken gconv.Struct error", zap.Error(err), zap.Any("token", token))
		return nil, *common.ServerErr.WithErr(err)
	}
	return adminUser, common.OK
}

// LarkBind 绑定飞书账号
func (s *Service) LarkBind(ctx context.Context, adminUser *common.AdminUser, req *dto.LarkBindReq) common.Errno {
	// 获取飞书用户 access token
	accessToken, errno := s.token.GetLarkUserAccessToken(ctx, req.AppCode, req.Code, req.RedirectUrl, "", false)
	if !errno.IsOK() {
		logger.Error("LarkQrCodeLogin GetLarkUserAccessToken error", zap.Any("req", req))
		return common.ServerErr
	}
	// 通过 access token 获取飞书用户信息
	larkUserInfo, err := s.lark.GetLarkUserInfo(ctx, accessToken.Token)
	if err != nil {
		logger.Error("LarkBind GetLarkUserInfo error", zap.Error(err), zap.Any("req", req))
		return *common.ServerErr.WithErr(err)
	}
	err = s.adminUser.UpdateUserLarkOpenID(ctx, adminUser.UserID, larkUserInfo.OpenID)
	if err != nil {
		logger.Error("LarkBind error", zap.Error(err), zap.Any("req", adminUser))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

// LarkUnbind 解绑飞书账号
func (s *Service) LarkUnbind(ctx context.Context, adminUser *common.AdminUser) common.Errno {
	err := s.adminUser.UpdateUserLarkOpenID(ctx, adminUser.UserID, "")
	if err != nil {
		logger.Error("LarkUnbind error", zap.Error(err), zap.Any("adminUser", adminUser))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}
