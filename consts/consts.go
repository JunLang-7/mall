package consts

import "time"

const (
	AdminTokenKey   = "admin_token"
	UserTokenKey    = "user_token"
	CustomerUserKey = "customer_user"
	AdminUserKey    = "admin_user"
)

const (
	IsEnable  = 1
	IsDisable = -1
)

const (
	ExpireLoginTime          = time.Minute * 2
	ExpireTicketTime         = time.Minute * 5
	ExpireAdminUserTokenTime = time.Hour * 24
	ExpirePasswordErrTime    = time.Minute * 5
	ExpireVerifyCodeErrTime  = time.Minute * 5
)

const (
	PasswordErrMaxCount = 3
)

const (
	WechatAppCode = 1000
	LarkAppCode   = 2000
)
