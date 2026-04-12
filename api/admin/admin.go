package admin

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/service/admin"
)

type Ctrl struct {
	adaptor adaptor.IAdaptor
	user    *admin.Service
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		user:    admin.NewService(adaptor),
	}
}
