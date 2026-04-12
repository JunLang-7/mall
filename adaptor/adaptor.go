package adaptor

import (
	"github.com/JunLang-7/mall/config"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type IAdaptor interface {
	GetConf() *config.Config
	GetDB() *gorm.DB
	GetRedis() *redis.Client
}

type Adaptor struct {
	conf *config.Config
	db   *gorm.DB
	rds  *redis.Client
}

func NewAdaptor(conf *config.Config, db *gorm.DB, rds *redis.Client) IAdaptor {
	return &Adaptor{
		conf: conf,
		db:   db,
		rds:  rds,
	}
}

func (a *Adaptor) GetConf() *config.Config {
	return a.conf
}

func (a *Adaptor) GetDB() *gorm.DB {
	return a.db
}

func (a *Adaptor) GetRedis() *redis.Client {
	return a.rds
}
