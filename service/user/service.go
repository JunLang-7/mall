package user

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/redis"
	"github.com/JunLang-7/mall/adaptor/rpc"
	"github.com/JunLang-7/mall/config"
	"github.com/JunLang-7/mall/service/token"
	"github.com/JunLang-7/mall/utils/captcha"
	"github.com/wenlng/go-captcha/v2/slide"
)

type Service struct {
	conf    *config.Config
	verify  redis.IVerify
	captcha slide.Captcha
	token   *token.Service
	lark    rpc.ILark
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		conf:    adaptor.GetConf(),
		verify:  redis.NewVerify(adaptor),
		captcha: captcha.NewSlideCaptcha(),
		token:   token.NewService(adaptor),
		lark:    rpc.NewLark(adaptor),
	}
}
