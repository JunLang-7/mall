package admin

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/service/admin"
	"github.com/JunLang-7/mall/service/perm"
	"github.com/JunLang-7/mall/service/role"
)

type Ctrl struct {
	adaptor adaptor.IAdaptor
	user    *admin.Service
	perm    *perm.Service
	role    *role.Service
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		user:    admin.NewService(adaptor),
		perm:    perm.NewService(adaptor),
		role:    role.NewService(adaptor),
	}
}
