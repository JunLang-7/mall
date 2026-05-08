package user

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/redis"
	"github.com/JunLang-7/mall/adaptor/repo/customer"
	"github.com/JunLang-7/mall/adaptor/rpc"
	"github.com/JunLang-7/mall/config"
	"github.com/JunLang-7/mall/service/token"
	"github.com/JunLang-7/mall/utils/captcha"
	"github.com/JunLang-7/mall/utils/snowflake"
	"github.com/wenlng/go-captcha/v2/slide"
)

type Service struct {
	conf      *config.Config
	verify    redis.IVerify
	captcha   slide.Captcha
	token     *token.Service
	lark      rpc.ILark
	storage   rpc.IStorage
	snow      *snowflake.Node
	orderFee  redis.IOrderFee
	qrcode    redis.IQrCode
	course    customer.ICourse
	lesson    customer.ILesson
	userRepo  customer.IUser
	order     customer.IOrder
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		conf:      adaptor.GetConf(),
		verify:    redis.NewVerify(adaptor),
		captcha:   captcha.NewSlideCaptcha(),
		token:     token.NewService(adaptor),
		lark:      rpc.NewLark(adaptor),
		storage:   rpc.NewStorage(adaptor),
		snow:      snowflake.NewNode(adaptor.GetConf().Order.SnowflakeNodeID),
		orderFee:  redis.NewOrderFee(adaptor),
		qrcode:    redis.NewQrCode(adaptor),
		course:    customer.NewCourseRepo(adaptor),
		lesson:    customer.NewLessonRepo(adaptor),
		userRepo:  customer.NewUserRepo(adaptor),
		order:     customer.NewOrderRepo(adaptor),
	}
}
