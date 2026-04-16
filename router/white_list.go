package router

// AdminAuthWhiteList 管理后台认证白名单，免登录访问
var AdminAuthWhiteList = map[string]bool{
	"/ping":                                true,
	"/metrics":                             true,
	"/admin/v1/user/verify/captcha/check":  true,
	"/admin/v1/user/verify/captcha":        true,
	"/admin/v1/user/verify/smscode":        true,
	"/admin/v1/user/mobile/verify_login":   true,
	"/admin/v1/user/mobile/password_login": true,
	"/admin/v1/user/lark/qrcode_login":     true,
	"/admin/v1/user/password/reset":        true,
}
