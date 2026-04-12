package admin

import (
	"context"

	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type IAdminUser interface {
	GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error)
}

type Repo struct {
	db  *gorm.DB
	rds *redis.Client
}

func NewRepo(db *gorm.DB, rds *redis.Client) *Repo {
	return &Repo{
		db:  db,
		rds: rds,
	}
}

func (r *Repo) GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error) {
	qs := query.Use(r.db).AdminUser
	return qs.WithContext(ctx).Where(qs.ID.Eq(userId)).First()
}
