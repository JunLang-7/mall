package admin

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/service/admin"
	"github.com/JunLang-7/mall/service/goods"
	"github.com/JunLang-7/mall/service/perm"
	"github.com/JunLang-7/mall/service/role"
	"github.com/JunLang-7/mall/service/storage"
	"github.com/JunLang-7/mall/service/user"
)

type Ctrl struct {
	adaptor adaptor.IAdaptor
	user    *admin.Service
	perm    *perm.Service
	role    *role.Service
	lesson  *goods.Service
	storage *storage.Service
	customer *user.Service
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		user:    admin.NewService(adaptor),
		perm:    perm.NewService(adaptor),
		role:    role.NewService(adaptor),
		lesson:  goods.NewService(adaptor),
		storage: storage.NewService(adaptor),
		customer: user.NewService(adaptor),
	}
}
