package perm

import (
	"context"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/admin"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/do"
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

// CreatePermission 创建权限
func (s *Service) CreatePermission(ctx context.Context, user *common.AdminUser, req *dto.PermissionCreateReq) (int64, common.Errno) {
	permID, err := s.adminPerm.CreatePermission(ctx, &do.PermCreate{
		AdminUserID: user.UserID,
		Code:        req.Code,
		Type:        req.Type,
		Name:        req.Name,
		PagePath:    req.PagePath,
		ParentID:    req.ParentID,
		Sort:        req.Sort,
		Desc:        req.Desc,
	})
	if err != nil {
		logger.Error("CreatePermission PermissionCreate error", zap.Error(err))
		return 0, *common.DataBaseErr.WithErr(err)
	}
	return permID, common.OK
}

// UpdatePermissions 批量更新权限
func (s *Service) UpdatePermissions(ctx context.Context, user *common.AdminUser, req *dto.PermissionUpdateResp) common.Errno {
	reqList := make([]do.PermUpdate, 0)
	lo.ForEach(req.List, func(item dto.PermissionUpdateDto, index int) {
		reqList = append(reqList, do.PermUpdate{
			ID: item.ID,
			PermCreate: do.PermCreate{
				AdminUserID: item.ID,
				Code:        item.Code,
				Type:        item.Type,
				Name:        item.Name,
				PagePath:    item.PagePath,
				ParentID:    lo.Ternary(item.ParentID == 0, -1, item.ParentID),
				Sort:        item.Sort,
				Desc:        item.Desc,
			},
		})
	})
	err := s.adminPerm.UpdatePermissions(ctx, &do.PermUpdateList{List: reqList})
	if err != nil {
		logger.Error("UpdatePermissions PermissionUpdateList error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

// DeletePermission 删除权限
func (s *Service) DeletePermission(ctx context.Context, user *common.AdminUser, req *dto.PermissionDeleteReq) common.Errno {
	err := s.adminPerm.DeletePermission(ctx, &do.PermDelete{ID: req.ID})
	if err != nil {
		logger.Error("DeletePermission PermissionDelete error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}
