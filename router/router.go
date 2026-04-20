package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/api/admin"
	"github.com/JunLang-7/mall/api/customer"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/config"
	"github.com/gin-gonic/gin"
)

type IRouter interface {
	Register(engine *gin.Engine)
	SpanFilter(r *gin.Context) bool
	AccessRecordFilter(r *gin.Context) bool
}

type Router struct {
	FullPPROF bool
	rootPath  string
	conf      *config.Config
	checkFunc func() error
	admin     *admin.Ctrl
	customer  *customer.Ctrl
}

func NewRouter(conf *config.Config, adaptor adaptor.IAdaptor, checkFunc func() error) *Router {
	return &Router{
		FullPPROF: conf.Server.EnablePprof,
		rootPath:  "/api/mall",
		conf:      conf,
		checkFunc: checkFunc,
		admin:     admin.NewCtrl(adaptor),
		customer:  customer.NewCtrl(adaptor),
	}
}

func (r *Router) Register(engine *gin.Engine) {
	if r.conf.Server.EnablePprof {
		SetupPprof(engine, "/debug/pprof")
	}
	engine.Any("/ping", r.checkServer())

	// 静态文件服务 (前端页面)
	engine.Static("/web", "./web")

	root := engine.Group(r.rootPath)
	r.route(root)
}

func (r *Router) SpanFilter(c *gin.Context) bool {
	path := strings.Replace(c.Request.URL.Path, r.rootPath, "", 1)
	if _, ok := AdminAuthWhiteList[path]; ok {
		return false
	}
	return true
}

func (r *Router) AccessRecordFilter(c *gin.Context) bool {
	return true
}

func (r *Router) route(root *gin.RouterGroup) {
	r.adminRoute(root)
	r.customerRoute(root)
}

func (r *Router) adminRoute(root *gin.RouterGroup) {
	adminRoot := root.Group("/admin", AdminAuthMiddleware(r.SpanFilter, func(ctx context.Context, token string) (*common.AdminUser, error) {
		return r.admin.GetAdminUserByToken(ctx, token)
	}))
	// 登录无鉴权 添加白名单
	adminRoot.GET("/v1/user/verify/captcha", r.admin.GetSmsCodeCaptcha)
	adminRoot.POST("/v1/user/verify/captcha/check", r.admin.CheckSmsCodeCaptcha)
	adminRoot.POST("/v1/user/verify/smscode", r.admin.GetSmsCodeVerify)
	adminRoot.POST("/v1/user/mobile/password_login", r.admin.MobilePasswordLogin)
	adminRoot.POST("/v1/user/mobile/verify_login", r.admin.MobileVerifyLogin)
	adminRoot.POST("/v1/user/lark/qrcode_login", r.admin.LarkQrCodeLogin)
	adminRoot.POST("/v1/user/mobile/reset_password", r.admin.MobilePasswordReset)

	// 管理员用户
	adminRoot.GET("/v1/user/info", r.admin.GetUserInfo)
	adminRoot.POST("/v1/user/create", r.admin.CreateUser)
	adminRoot.POST("/v1/user/update", r.admin.UpdateUser)
	adminRoot.POST("/v1/user/update-status", r.admin.UpdateUserStatus)
	adminRoot.POST("/v1/user/delete", r.admin.DeleteUser)
	adminRoot.POST("/v1/user/logout", r.admin.AdminUserLogout)

	// 绑定解绑飞书账号、修改手机号等敏感操作
	adminRoot.POST("/v1/user/lark_bind", r.admin.LarkBind)
	adminRoot.POST("/v1/user/lark_unbind", r.admin.LarkUnbind)

	// 权限菜单
	adminRoot.GET("/v1/perm/list", r.admin.PermissionList)
	adminRoot.GET("/v1/perm/my_perms", r.admin.MyPermissionList)
	adminRoot.POST("/v1/perm/create", r.admin.CreatePermission)
	adminRoot.POST("/v1/perm/update", r.admin.UpdatePermission)
	adminRoot.POST("/v1/perm/delete", r.admin.DeletePermission)
}

func (r *Router) customerRoute(root *gin.RouterGroup) {
	cstRoot := root.Group("/customer", AuthMiddleware(r.SpanFilter, func(ctx context.Context, token string) (*common.User, error) {
		return &common.User{}, nil
	}))
	cstRoot.Any("/user/info", r.admin.GetUserInfo)
}

func (r *Router) checkServer() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		err := r.checkFunc()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}
