package admin

import (
	"context"
	"errors"

	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *Service) GetUserInfo(ctx context.Context, adminUser *common.AdminUser) (*dto.UserInfoResp, common.Errno) {
	user, err := s.adminUser.GetUserInfo(ctx, 1)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.UserNotFoundErr
		}
		logger.Error("GetUserInfo error", zap.Error(err), zap.Any("user_id", adminUser.UserID))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return &dto.UserInfoResp{UserID: user.ID, Name: user.Name}, common.OK
}
