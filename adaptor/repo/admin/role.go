package admin

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/utils/tools"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type IRole interface {
	CreateRole(ctx context.Context, req *do.AddRole) (int64, error)
	UpdateRole(ctx context.Context, req *do.UpdateRole) error
	SetRolePerms(ctx context.Context, roleID int64, permIDs []int64, userID int64) error
	GetRolePerms(ctx context.Context, roleIDs []int64) (map[int64][]int64, error)
	ListRoles(ctx context.Context, req *do.ListRole) ([]*model.Role, int64, error)
	GetRoleByUserID(ctx context.Context, userId int64) ([]*model.AdminUserRole, error)
	GetRoleByUserIDs(ctx context.Context, userIDs []int64) (map[int64][]*model.AdminUserRole, error)
	GetRoleByIDs(ctx context.Context, roleIDs []int64) (map[int64]*model.Role, error)
}

type Role struct {
	db *gorm.DB
}

func NewAdminRole(adaptor adaptor.IAdaptor) *Role {
	return &Role{
		db: adaptor.GetDB(),
	}
}

func (r *Role) CreateRole(ctx context.Context, req *do.AddRole) (int64, error) {
	timeNow := time.Now()
	qs := query.Use(r.db).Role
	addObj := &model.Role{
		Name:     req.Name,
		Desc:     req.Desc,
		Status:   consts.IsEnable,
		CreateAt: timeNow,
		UpdateAt: timeNow,
		CreateBy: req.AdminUserID,
		UpdateBy: req.AdminUserID,
	}
	err := qs.WithContext(ctx).Create(addObj)
	if err != nil {
		return 0, err
	}
	return addObj.ID, nil
}

func (r *Role) UpdateRole(ctx context.Context, req *do.UpdateRole) error {
	qs := query.Use(r.db).Role
	updateMap := map[string]interface{}{
		qs.Name.ColumnName().String():     req.Name,
		qs.Desc.ColumnName().String():     req.Desc,
		qs.Status.ColumnName().String():   req.Status,
		qs.UpdateAt.ColumnName().String(): time.Now(),
		qs.UpdateBy.ColumnName().String(): req.AdminUserID,
	}
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(updateMap)
	if err != nil {
		return err
	}
	return nil
}

func (r *Role) SetRolePerms(ctx context.Context, roleID int64, permIDs []int64, userID int64) error {
	timeNow := time.Now()
	qs := query.Use(r.db).RolePermission
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.RolePermission{}).Where(qs.RoleID.Eq(roleID)).Delete(&model.RolePermission{}).Error
		if err != nil {
			return err
		}
		rolePerms := make([]*model.RolePermission, 0)
		for _, permID := range permIDs {
			rolePerms = append(rolePerms, &model.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
				CreateAt:     timeNow,
				UpdateAt:     timeNow,
				CreateBy:     userID,
				UpdateBy:     userID,
			})
		}
		return tx.CreateInBatches(rolePerms, len(rolePerms)).Error
	})
}

func (r *Role) GetRolePerms(ctx context.Context, roleIDs []int64) (map[int64][]int64, error) {
	qs := query.Use(r.db).RolePermission
	list, err := qs.WithContext(ctx).Where(qs.RoleID.In(roleIDs...)).Find()
	if err != nil {
		return nil, err
	}
	rolePermMaps := make(map[int64][]int64)
	lo.ForEach(list, func(item *model.RolePermission, index int) {
		rolePermMaps[item.RoleID] = append(rolePermMaps[item.RoleID], item.PermissionID)
	})
	return rolePermMaps, nil
}

func (r *Role) ListRoles(ctx context.Context, req *do.ListRole) ([]*model.Role, int64, error) {
	qs := query.Use(r.db).Role
	tx := qs.WithContext(ctx)
	if req.NameKw != "" {
		tx = tx.Where(qs.Name.Like(tools.GetAllLike(req.NameKw)))
	}
	if req.Status != 0 {
		tx = tx.Where(qs.Status.Eq(req.Status))
	}
	return tx.Order(qs.Status.Desc(), qs.CreateAt.Desc()).FindByPage(req.GetOffset(), req.Limit)
}

func (r *Role) GetRoleByUserID(ctx context.Context, userId int64) ([]*model.AdminUserRole, error) {
	qs := query.Use(r.db).AdminUserRole
	return qs.WithContext(ctx).Where(qs.AdminUserID.Eq(userId)).Find()
}

func (r *Role) GetRoleByUserIDs(ctx context.Context, userIDs []int64) (map[int64][]*model.AdminUserRole, error) {
	qs := query.Use(r.db).AdminUserRole
	list, err := qs.WithContext(ctx).Where(qs.AdminUserID.In(userIDs...)).Find()
	if err != nil {
		return nil, err
	}
	roleMap := lo.GroupBy(list, func(item *model.AdminUserRole) int64 {
		return item.AdminUserID
	})
	return roleMap, nil
}

func (r *Role) GetRoleByIDs(ctx context.Context, roleIDs []int64) (map[int64]*model.Role, error) {
	qs := query.Use(r.db).Role
	list, err := qs.WithContext(ctx).Where(qs.ID.In(roleIDs...)).Find()
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(list, func(item *model.Role) (int64, *model.Role) {
		return item.ID, item
	}), nil
}
