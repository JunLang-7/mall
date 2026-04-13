package admin

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/redis"
	"github.com/JunLang-7/mall/adaptor/repo/admin"
	"github.com/JunLang-7/mall/utils/captcha"
	"github.com/wenlng/go-captcha/v2/slide"
)

type Service struct {
	adminUser admin.IAdminUser
	verify    redis.IVerify
	captcha   slide.Captcha
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminUser: admin.NewRepo(adaptor.GetDB(), adaptor.GetRedis()),
		verify:    redis.NewVerify(adaptor),
		captcha:   captcha.NewSlideCaptcha(),
	}
}
