package perm

import (
	"context"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/admin"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Service struct {
	adminPerm admin.IPerm
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminPerm: admin.NewAdminPerm(adaptor),
	}
}

// PermissionList 获取权限列表
func (s *Service) PermissionList(ctx context.Context) (*dto.PermissionListResp, common.Errno) {
	permList, total, err := s.adminPerm.PermissionList(ctx, common.Pager{
		Page:      1,
		Limit:     1000,
		Unlimited: true,
	})
	if err != nil {
		logger.Error("PermissionList PermissionList error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	retList := make([]*dto.PermissionDto, 0, len(permList))
	lo.ForEach(permList, func(item *model.Permission, index int) {
		retList = append(retList, &dto.PermissionDto{
			ID:       item.ID,
			Code:     item.Code,
			Type:     item.Type,
			Name:     item.Name,
			PagePath: item.PagePath,
			ParentID: item.ParentID,
			Status:   item.Status,
			Sort:     item.Sort,
			Desc:     item.Desc,
		})
	})
	return &dto.PermissionListResp{
		Pager: common.Pager{
			Page:      1,
			Limit:     1000,
			Unlimited: true,
		},
		Total: total,
		List:  retList,
	}, common.OK
}

// MyPermissionList 获取我的权限列表
func (s *Service) MyPermissionList(ctx context.Context, user *common.AdminUser) ([]*dto.PermissionDto, common.Errno) {
	permList, err := s.adminPerm.MyPermissionList(ctx, user.UserID)
	if err != nil {
		logger.Error("MyPermissionList PermissionList error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	retList := make([]*dto.PermissionDto, 0, len(permList))
	lo.ForEach(permList, func(item *model.Permission, index int) {
		retList = append(retList, &dto.PermissionDto{
			ID:       item.ID,
			Code:     item.Code,
			Type:     item.Type,
			Name:     item.Name,
			PagePath: item.PagePath,
			ParentID: item.ParentID,
			Status:   item.Status,
			Sort:     item.Sort,
			Desc:     item.Desc,
		})
	})
	return retList, common.OK
}
