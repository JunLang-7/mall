package consts

import "time"

const (
	AdminTokenKey   = "token"
	UserTokenKey    = "token"
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
	ExpireCustomerTokenTime  = time.Hour * 24 * 7
	ExpirePasswordErrTime    = time.Minute * 5
	ExpireVerifyCodeErrTime  = time.Minute * 5
	ExpireOrderFeeTime       = time.Minute * 10
	ExpireOrderPayTime       = time.Minute * 30
)

const (
	PasswordErrMaxCount = 3
)

const (
	WechatAppCode = 1000
	LarkAppCode   = 2000
)

// 验证码场景
const (
	AddAdminUserPasswordSmsCode   = "add_admin_user_password"
	AdminUserMobileLoginSmsCode   = "admin_user_mobile_login"
	AdminUserResetPasswordSmsCode = "admin_user_reset_password"
	AdminUserChangeMobileSmsCode  = "admin_user_change_mobile"

	CustomerMobileLoginSmsCode     = "customer_mobile_login"
	CustomerResetPasswordSmsCode   = "customer_reset_password"
	CustomerChangePasswordSmsCode  = "customer_change_password"
	CustomerRegisterLoginSmsCode   = "customer_mobile_login"
)

const (
	OrderStatusCanceled = -1
	OrderStatusPending  = 1
	OrderStatusPaid     = 2
	OrderStatusRefunded = 3
	OrderStatusShipped  = 4
	OrderStatusSigned   = 5
	OrderStatusDone     = 6
)

const (
	OrderSourceCustomer = 1
	OrderSourceAdmin    = 2
	OrderSourceSystem   = 3
)

const (
	GoodsTypeCourse = 1
)

const (
	CancelTypeUser    = 1
	CancelTypeAdmin   = 2
	CancelTypeTimeout = 3
	SystemUserID       = -1
)

const (
	RefundStatusNone    = 0
	RefundStatusPending = 1
	RefundStatusDone    = 2
	RefundStatusError   = 3
)

const (
	ReceiveConfirmUser = 1
	ReceiveConfirmAuto = 99
)

const (
	LearnStatusLearning = 1
	LearnStatusDone     = 2
)
