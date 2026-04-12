package admin

import (
	"context"

	"github.com/JunLang-7/mall/service/do"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type IAdminUser interface {
	Hello(ctx context.Context, req *do.Hello) (string, error)
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

func (r *Repo) Hello(ctx context.Context, req *do.Hello) (res string, err error) {
	return "hello world", nil
}
