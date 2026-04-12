package admin

import (
	"context"
	"errors"

	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *Service) GetUserInfo(ctx context.Context, adminUser *common.AdminUser, userID int64) (*dto.UserInfoResp, common.Errno) {
	user, err := s.adminUser.GetUserInfo(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.UserNotFoundErr
		}
		logger.Error("GetUserInfo error", zap.Error(err), zap.Int64("request_user_id", userID), zap.Any("admin_user", adminUser))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return &dto.UserInfoResp{UserID: user.ID, Name: user.Name}, common.OK
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
