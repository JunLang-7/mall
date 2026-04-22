package admin

import (
	"context"
	"errors"

	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/go-redis/redis"
	"github.com/gogf/gf/util/gconv"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *Service) GetUserInfo(ctx context.Context, adminUser *common.AdminUser, userID int64) (*dto.AdminUserWithRoleDto, common.Errno) {
	user, err := s.adminUser.GetUserInfo(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.InvalidPasswordErr
		}
		logger.Error("GetUserInfo error", zap.Error(err), zap.Int64("request_user_id", userID), zap.Any("admin_user", adminUser))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	userRoleMap, err := s.adminRole.GetRoleByUserIDs(ctx, []int64{user.ID})
	if err != nil {
		logger.Error("GetUserInfo GetRoleByUserID error", zap.Error(err), zap.Any("userID", userID))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	roleIDs := make([]int64, 0)
	for _, vList := range userRoleMap {
		for _, v := range vList {
			roleIDs = append(roleIDs, v.RoleID)
		}
	}
	roleMap, err := s.adminRole.GetRoleByIDs(ctx, lo.Uniq(roleIDs))
	if err != nil {
		logger.Error("GetUserInfo GetRoleByUserIDs error", zap.Error(err), zap.Any("user", adminUser))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	roles := make([]*common.IDName, 0)
	for _, roleID := range roleIDs {
		roles = append(roles, &common.IDName{
			ID:   roleMap[roleID].ID,
			Name: roleMap[roleID].Name,
		})
	}
	return &dto.AdminUserWithRoleDto{
		AdminUserDto: dto.AdminUserDto{
			UserID:     user.ID,
			Name:       user.Name,
			NickName:   user.NickName,
			Sex:        user.Sex,
			Status:     user.Status,
			Mobile:     user.Mobile,
			LarkOpenID: user.LarkOpenID,
			UpdateAt:   user.UpdateAt.UnixMilli(),
			CreateAt:   user.CreateAt.UnixMilli(),
		},
		Roles: roles,
	}, common.OK
}

func (s *Service) ListUsers(ctx context.Context, adminUser *common.AdminUser, req *dto.ListUsersReq) (*dto.ListUsersResp, common.Errno) {
	userList, total, err := s.adminUser.ListUsers(ctx, &do.ListUsers{
		Name:   req.Name,
		Mobile: req.Mobile,
		RoleID: req.RoleID,
		Status: req.Status,
		Pager:  req.Pager,
	})
	if err != nil {
		logger.Error("ListUsers error", zap.Error(err), zap.Any("admin_user", adminUser), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	userIDs := make([]int64, 0)
	lo.ForEach(userList, func(item *model.AdminUser, index int) {
		userIDs = append(userIDs, item.ID)
	})

	userRoleMap, err := s.adminRole.GetRoleByUserIDs(ctx, userIDs)
	if err != nil {
		logger.Error("ListUsers GetRoleByUserIDs error", zap.Error(err), zap.Any("admin_user", adminUser), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	roleIDs := make([]int64, 0)
	for _, vList := range userRoleMap {
		for _, v := range vList {
			roleIDs = append(roleIDs, v.RoleID)
		}
	}
	roleMap, err := s.adminRole.GetRoleByIDs(ctx, lo.Uniq(roleIDs))
	if err != nil {
		logger.Error("ListUsers GetRoleByIDs error", zap.Error(err), zap.Any("admin_user", adminUser), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	retList := make([]*dto.AdminUserWithRoleDto, 0, len(userList))
	lo.ForEach(userList, func(user *model.AdminUser, index int) {
		retList = append(retList, &dto.AdminUserWithRoleDto{
			AdminUserDto: dto.AdminUserDto{
				UserID:     user.ID,
				Name:       user.Name,
				NickName:   user.NickName,
				Sex:        user.Sex,
				Status:     user.Status,
				Mobile:     user.Mobile,
				LarkOpenID: user.LarkOpenID,
				UpdateAt:   user.UpdateAt.UnixMilli(),
				CreateAt:   user.CreateAt.UnixMilli(),
			},
			Roles: lo.Map(userRoleMap[user.ID], func(item *model.AdminUserRole, index int) *common.IDName {
				return &common.IDName{
					ID:   roleMap[item.RoleID].ID,
					Name: roleMap[item.RoleID].Name,
				}
			}),
		})
	})
	return &dto.ListUsersResp{
		Pager: req.Pager,
		Total: total,
		List:  retList,
	}, common.OK
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

// AdminUserLogout 管理员用户登出
func (s *Service) AdminUserLogout(ctx context.Context, adminUser *common.AdminUser) common.Errno {
	err := s.verify.CleanToken(ctx, adminUser.UserID)
	if err != nil {
		logger.Error("AdminUserLogout error", zap.Error(err), zap.Any("adminUser", adminUser))
		return *common.RedisErr.WithErr(err)
	}
	return common.OK
}
