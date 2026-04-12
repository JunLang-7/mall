package router

import (
	"net/http"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/api/admin"
	"github.com/JunLang-7/mall/api/customer"
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
	engine.Use(AuthMiddleware(r.SpanFilter))
	if r.conf.Server.EnablePprof {
		SetupPprof(engine, "/debug/pprof")
	}
	engine.Any("/ping", r.checkServer())

	root := engine.Group(r.rootPath)
	r.route(root)
}

func (r *Router) SpanFilter(c *gin.Context) bool {
	return true
}

func (r *Router) AccessRecordFilter(c *gin.Context) bool {
	return true
}

func (r *Router) route(root *gin.RouterGroup) {
	adminRoot := root.Group("/admin")
	adminRoot.GET("/user/info", r.admin.GetUserInfo)
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
