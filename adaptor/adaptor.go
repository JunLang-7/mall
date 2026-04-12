package adaptor

import (
	"github.com/JunLang-7/mall/config"
	"gorm.io/gorm"
)

type IAdaptor interface {
	GetDB() *gorm.DB
}

type Adaptor struct {
	conf *config.Config
	db   *gorm.DB
}

func NewAdaptor(conf *config.Config, db *gorm.DB) IAdaptor {
	return &Adaptor{
		conf: conf,
		db:   db,
	}
}

func (a *Adaptor) GetConf() *config.Config {
	return a.conf
}

func (a *Adaptor) GetDB() *gorm.DB {
	return a.db
}
