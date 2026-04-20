package admin

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"gorm.io/gorm"
)

type IPerm interface {
	PermissionList(ctx context.Context, pager common.Pager) ([]*model.Permission, int64, error)
	MyPermissionList(ctx context.Context, adminUserId int64) ([]*model.Permission, error)
	CreatePermission(ctx context.Context, req *do.PermCreate) (int64, error)
	UpdatePermissions(ctx context.Context, req *do.PermUpdateList) error
	DeletePermission(ctx context.Context, req *do.PermDelete) error
}

type AdminPerm struct {
	db *gorm.DB
}

func NewAdminPerm(adaptor adaptor.IAdaptor) *AdminPerm {
	return &AdminPerm{
		db: adaptor.GetDB(),
	}
}

func (a *AdminPerm) PermissionList(ctx context.Context, pager common.Pager) ([]*model.Permission, int64, error) {
	qs := query.Use(a.db).Permission
	list, total, err := qs.WithContext(ctx).Where(qs.Status.Eq(consts.IsEnable)).FindByPage(pager.GetOffset(), pager.Limit)
	return list, total, err
}

func (a *AdminPerm) MyPermissionList(ctx context.Context, adminUserId int64) ([]*model.Permission, error) {
	qs := query.Use(a.db).Permission
	return qs.WithContext(ctx).Where(qs.Status.Eq(consts.IsEnable)).Where(qs.ParentID.Eq(adminUserId)).Find()
}

func (a *AdminPerm) CreatePermission(ctx context.Context, req *do.PermCreate) (int64, error) {
	qs := query.Use(a.db).Permission
	timeNow := time.Now()
	perm := &model.Permission{
		Code:     req.Code,
		Type:     req.Type,
		Name:     req.Name,
		PagePath: req.PagePath,
		ParentID: req.ParentID,
		Status:   consts.IsEnable,
		Sort:     req.Sort,
		Desc:     req.Desc,
		CreateAt: timeNow,
		UpdateAt: timeNow,
		UpdateBy: req.AdminUserID,
	}
	err := qs.WithContext(ctx).Create(perm)
	if err != nil {
		return 0, err
	}
	return perm.ID, nil
}

func (a *AdminPerm) UpdatePermissions(ctx context.Context, req *do.PermUpdateList) error {
	qs := query.Use(a.db).Permission
	return a.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range req.List {
			updateMap := map[string]interface{}{
				qs.UpdateBy.ColumnName().String(): item.AdminUserID,
				qs.UpdateAt.ColumnName().String(): time.Now(),
			}
			if item.Code != "" {
				updateMap[qs.Code.ColumnName().String()] = item.Code
			}
			if item.Type != 0 {
				updateMap[qs.Type.ColumnName().String()] = item.Type
			}
			if item.Name != "" {
				updateMap[qs.Name.ColumnName().String()] = item.Name
			}
			if item.PagePath != "" {
				updateMap[qs.PagePath.ColumnName().String()] = item.PagePath
			}
			if item.ParentID != 0 {
				updateMap[qs.ParentID.ColumnName().String()] = item.ParentID
			}
			if item.Sort != 0 {
				updateMap[qs.Sort.ColumnName().String()] = item.Sort
			}
			if item.Desc != "" {
				updateMap[qs.Desc.ColumnName().String()] = item.Desc
			}
			err := tx.Model(&model.Permission{}).Where(qs.ID.Eq(item.ID)).Updates(updateMap).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *AdminPerm) DeletePermission(ctx context.Context, req *do.PermDelete) error {
	qs := query.Use(a.db).Permission
	_, err := qs.WithContext(ctx).Delete(&model.Permission{
		ID: req.ID,
	})
	return err
}
