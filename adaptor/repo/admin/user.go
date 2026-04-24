package admin

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"github.com/go-redis/redis"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type IAdminUser interface {
	GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error)
	ListUsers(ctx context.Context, req *do.ListUsers) ([]*model.AdminUser, int64, error)
	GetUserByMobile(ctx context.Context, mobile string) (*model.AdminUser, error)
	GetUserByLarkOpenID(ctx context.Context, openID string) (*model.AdminUser, error)
	CreateUser(ctx context.Context, user *do.CreateUser) (int64, error)
	UpdateUser(ctx context.Context, user *do.UpdateUser) error
	DeleteUser(ctx context.Context, userId int64) error
	UpdateUserLarkOpenID(ctx context.Context, userId int64, openID string) error
	UpdateUserPassword(ctx context.Context, req *do.UpdateUserPassword) error
	GetUserNameMap(ctx context.Context, ids []int64) (map[int64]string, error)
}

type Repo struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewRepo(adaptor adaptor.IAdaptor) *Repo {
	return &Repo{
		db:  adaptor.GetDB(),
		rds: adaptor.GetRedis(),
	}
}

func (r *Repo) GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error) {
	qs := query.Use(r.db).AdminUser
	return qs.WithContext(ctx).Where(qs.ID.Eq(userId)).First()
}

func (r *Repo) ListUsers(ctx context.Context, req *do.ListUsers) ([]*model.AdminUser, int64, error) {
	qs := query.Use(r.db).AdminUser
	lists, total, err := qs.WithContext(ctx).Order(qs.CreateAt.Asc()).FindByPage(req.Pager.GetOffset(), req.Pager.Limit)
	return lists, total, err
}

func (r *Repo) GetUserByMobile(ctx context.Context, mobile string) (*model.AdminUser, error) {
	qs := query.Use(r.db).AdminUser
	return qs.WithContext(ctx).Where(qs.Mobile.Eq(mobile)).First()
}

func (r *Repo) GetUserByLarkOpenID(ctx context.Context, openID string) (*model.AdminUser, error) {
	qs := query.Use(r.db).AdminUser
	return qs.WithContext(ctx).Where(qs.LarkOpenID.Eq(openID)).First()
}

func (r *Repo) CreateUser(ctx context.Context, req *do.CreateUser) (int64, error) {
	timeNow := time.Now()
	qs := query.Use(r.db).AdminUser
	addObj := &model.AdminUser{
		Name:     req.Name,
		NickName: req.NickName,
		Mobile:   req.Mobile,
		Status:   consts.IsEnable,
		Sex:      req.Sex,
		CreateAt: timeNow,
		UpdateAt: timeNow,
		CreateBy: req.AdminUserID,
		UpdateBy: req.AdminUserID,
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := qs.WithContext(ctx).Create(addObj)
		if err != nil {
			return err
		}
		userRoles := make([]model.AdminUserRole, 0)
		for _, roleID := range req.RoleIDs {
			userRoles = append(userRoles, model.AdminUserRole{
				AdminUserID: addObj.ID,
				RoleID:      roleID,
				UpdateAt:    timeNow,
				UpdateBy:    req.AdminUserID,
			})
		}
		return tx.CreateInBatches(userRoles, 100).Error
	})
	if err != nil {
		return 0, err
	}
	return addObj.ID, nil
}

func (r *Repo) UpdateUser(ctx context.Context, req *do.UpdateUser) error {
	qs := query.Use(r.db).AdminUser
	rqs := query.Use(r.db).AdminUserRole
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where(qs.ID.Eq(req.ID)).Updates(model.AdminUser{
			Name:     req.Name,
			NickName: req.NickName,
			UpdateAt: time.Now(),
			UpdateBy: req.AdminUserID,
			Sex:      req.Sex,
			Status:   req.Status,
		}).Error
		if err != nil {
			return err
		}
		err = tx.Where(rqs.AdminUserID.Eq(req.ID)).Delete(model.AdminUserRole{}).Error
		if err != nil {
			return err
		}
		userRoles := make([]model.AdminUserRole, 0)
		for _, roleID := range req.RoleIDs {
			userRoles = append(userRoles, model.AdminUserRole{
				AdminUserID: req.AdminUserID,
				RoleID:      roleID,
				UpdateAt:    time.Now(),
				UpdateBy:    req.AdminUserID,
			})
		}
		return tx.CreateInBatches(userRoles, 100).Error
	})
}

func (r *Repo) DeleteUser(ctx context.Context, userId int64) error {
	qs := query.Use(r.db).AdminUser
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(userId)).Delete()
	return err
}

func (r *Repo) UpdateUserLarkOpenID(ctx context.Context, userId int64, openID string) error {
	qs := query.Use(r.db).AdminUser
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(userId)).Update(qs.LarkOpenID, openID)
	return err
}

func (r *Repo) UpdateUserPassword(ctx context.Context, req *do.UpdateUserPassword) error {
	qs := query.Use(r.db).AdminUser
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Update(qs.Password, req.Password)
	return err
}

func (r *Repo) GetUserNameMap(ctx context.Context, ids []int64) (map[int64]string, error) {
	qs := query.Use(r.db).AdminUser
	list, err := qs.WithContext(ctx).Select(qs.ID, qs.Name).Where(qs.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(list, func(item *model.AdminUser) (int64, string) {
		return item.ID, item.Name
	}), nil
}
