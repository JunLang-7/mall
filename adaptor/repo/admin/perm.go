package admin

import (
	"context"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/consts"
	"gorm.io/gorm"
)

type IPerm interface {
	PermissionList(ctx context.Context, pager common.Pager) ([]*model.Permission, int64, error)
	MyPermissionList(ctx context.Context, adminUserId int64) ([]*model.Permission, error)
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
