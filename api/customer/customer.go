package customer

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/service/user"
)

type Ctrl struct {
	adaptor adaptor.IAdaptor
	user    *user.Service
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		user:    user.NewService(adaptor),
	}
}
