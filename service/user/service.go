package user

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/redis"
	"github.com/JunLang-7/mall/adaptor/rpc"
	"github.com/JunLang-7/mall/config"
	"github.com/JunLang-7/mall/service/token"
	"github.com/JunLang-7/mall/utils/captcha"
	"github.com/JunLang-7/mall/utils/snowflake"
	goredis "github.com/go-redis/redis"
	"github.com/wenlng/go-captcha/v2/slide"
	"gorm.io/gorm"
)

type Service struct {
	conf    *config.Config
	db      *gorm.DB
	rds     *goredis.Client
	verify  redis.IVerify
	captcha slide.Captcha
	token   *token.Service
	lark    rpc.ILark
	storage rpc.IStorage
	snow    *snowflake.Node
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		conf:    adaptor.GetConf(),
		db:      adaptor.GetDB(),
		rds:     adaptor.GetRedis(),
		verify:  redis.NewVerify(adaptor),
		captcha: captcha.NewSlideCaptcha(),
		token:   token.NewService(adaptor),
		lark:    rpc.NewLark(adaptor),
		storage: rpc.NewStorage(adaptor),
		snow:    snowflake.NewNode(adaptor.GetConf().Order.SnowflakeNodeID),
	}
}
