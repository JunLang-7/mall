package customer

import "github.com/JunLang-7/mall/adaptor"

type Ctrl struct {
	adaptor adaptor.IAdaptor
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{adaptor: adaptor}
}
