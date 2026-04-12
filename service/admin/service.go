package admin

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/admin"
)

type Service struct {
	adminUser admin.IAdminUser
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminUser: admin.NewRepo(adaptor.GetDB(), adaptor.GetRedis()),
	}
}
