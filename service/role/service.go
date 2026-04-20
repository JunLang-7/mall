package role

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
	adminRole admin.IRole
	adminPerm admin.IPerm
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminRole: admin.NewAdminRole(adaptor),
		adminPerm: admin.NewAdminPerm(adaptor),
	}
}

// ListRole 角色列表
func (s *Service) ListRole(ctx context.Context, req *dto.ListRoleReq) (*dto.ListRoleResp, common.Errno) {
	list, total, err := s.adminRole.ListRoles(ctx, &do.ListRole{
		Pager:  req.Pager,
		NameKw: req.NameKw,
		Status: req.Status,
	})
	if err != nil {
		logger.Error("ListRoles error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	roleIDs := make([]int64, 0)
	lo.ForEach(list, func(item *model.Role, index int) {
		roleIDs = append(roleIDs, item.ID)
	})

	rolePermMap, err := s.adminRole.GetRolePerms(ctx, lo.Uniq(roleIDs))
	if err != nil {
		logger.Error("ListRole GetRolePerms error", zap.Error(err), zap.Int64s("roleIDs", roleIDs))
		return nil, *common.DataBaseErr.WithErr(err)
	}

	permIDs := make([]int64, 0)
	for _, vList := range rolePermMap {
		permIDs = append(permIDs, vList...)
	}
	permNameMap, err := s.adminPerm.GetPermNameMap(ctx, permIDs)
	if err != nil {
		logger.Error("ListRole GetPermName error", zap.Error(err), zap.Int64s("permIDs", permIDs))
		return nil, *common.DataBaseErr.WithErr(err)
	}

	retList := make([]*dto.RoleDto, 0, len(list))
	lo.ForEach(list, func(item *model.Role, index int) {
		perms := make([]common.IDName, 0)
		lo.ForEach(rolePermMap[item.ID], func(item int64, index int) {
			perms = append(perms, common.IDName{
				ID:   item,
				Name: permNameMap[item],
			})
		})
		retList = append(retList, &dto.RoleDto{
			ID:       item.ID,
			Name:     item.Name,
			Desc:     item.Desc,
			Status:   item.Status,
			Perms:    perms,
			CreateAt: item.CreateAt.UnixMilli(),
			UpdateAt: item.UpdateAt.UnixMilli(),
		})
	})
	return &dto.ListRoleResp{
		Pager: req.Pager,
		Total: total,
		List:  retList,
	}, common.OK
}

// GetMyRoles 我的角色列表
func (s *Service) GetMyRoles(ctx context.Context, user *common.AdminUser) ([]*dto.RoleDto, common.Errno) {
	roleList, err := s.adminRole.GetRoleByUserID(ctx, user.UserID)
	if err != nil {
		logger.Error("GetMyRoles GetRoleByUserID error", zap.Error(err), zap.Any("userID", user.UserID))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	roleIDs := make([]int64, 0)
	lo.ForEach(roleList, func(item *model.AdminUserRole, index int) {
		roleIDs = append(roleIDs, item.RoleID)
	})
	roleIDs = lo.Uniq(roleIDs)

	roleMap, err := s.adminRole.GetRoleByIDs(ctx, roleIDs)
	if err != nil {
		logger.Error("GetMyRoles GetRoleByIDs error", zap.Error(err))
		return nil, *common.DataBaseErr.WithErr(err)
	}

	rolePermMap, err := s.adminRole.GetRolePerms(ctx, roleIDs)
	if err != nil {
		logger.Error("ListRole GetRolePerms error", zap.Error(err), zap.Int64s("roleIDs", roleIDs))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	permIDs := make([]int64, 0)
	for _, vList := range rolePermMap {
		permIDs = append(permIDs, vList...)
	}
	permNameMap, err := s.adminPerm.GetPermNameMap(ctx, permIDs)
	if err != nil {
		logger.Error("ListRole GetPermName error", zap.Error(err), zap.Int64s("permIDs", permIDs))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	retList := make([]*dto.RoleDto, 0, len(roleList))
	lo.ForEach(roleList, func(item *model.AdminUserRole, index int) {
		role, ok := roleMap[item.RoleID]
		if !ok {
			role = &model.Role{}
		}
		retList = append(retList, &dto.RoleDto{
			ID:     item.ID,
			Name:   role.Name,
			Desc:   role.Desc,
			Status: role.Status,
			Perms: lo.Map(rolePermMap[item.ID], func(item int64, index int) common.IDName {
				return common.IDName{
					ID:   item,
					Name: permNameMap[item],
				}
			}),
			CreateAt: role.CreateAt.UnixMilli(),
			UpdateAt: role.UpdateAt.UnixMilli(),
		})
	})
	return retList, common.OK
}

// CreateRole 创建角色
func (s *Service) CreateRole(ctx context.Context, user *common.AdminUser, req *dto.AddRoleReq) (int64, common.Errno) {
	permID, err := s.adminRole.CreateRole(ctx, &do.AddRole{
		AdminUserID: user.UserID,
		Name:        req.Name,
		Desc:        req.Desc,
	})
	if err != nil {
		logger.Error("CreateRole error", zap.Error(err))
		return 0, *common.DataBaseErr.WithErr(err)
	}
	return permID, common.OK
}

// UpdateRole 更新角色
func (s *Service) UpdateRole(ctx context.Context, user *common.AdminUser, req *dto.UpdateRoleReq) common.Errno {
	err := s.adminRole.UpdateRole(ctx, &do.UpdateRole{
		AdminUserID: user.UserID,
		ID:          req.ID,
		Name:        req.Name,
		Desc:        req.Desc,
		Status:      req.Status,
	})
	if err != nil {
		logger.Error("UpdateRoles error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

// SetRolePerms 设置角色权限
func (s *Service) SetRolePerms(ctx context.Context, user *common.AdminUser, req *dto.SetRolePermReq) common.Errno {
	err := s.adminRole.SetRolePerms(ctx, req.RoleID, req.PermIDs, user.UserID)
	if err != nil {
		logger.Error("SetRolePerms error", zap.Error(err))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}
