package admin

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/redis"
	"github.com/JunLang-7/mall/adaptor/repo/admin"
	"github.com/JunLang-7/mall/adaptor/rpc"
	"github.com/JunLang-7/mall/config"
	"github.com/JunLang-7/mall/service/token"
	"github.com/JunLang-7/mall/utils/captcha"
	"github.com/wenlng/go-captcha/v2/slide"
)

type Service struct {
	conf      *config.Config
	adminUser admin.IAdminUser
	verify    redis.IVerify
	captcha   slide.Captcha
	token     *token.Service
	lark      rpc.ILark
	adminRole admin.IRole
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		conf:      adaptor.GetConf(),
		adminUser: admin.NewRepo(adaptor),
		verify:    redis.NewVerify(adaptor),
		captcha:   captcha.NewSlideCaptcha(),
		token:     token.NewService(adaptor),
		lark:      rpc.NewLark(adaptor),
		adminRole: admin.NewAdminRole(adaptor),
	}
}
