package router

import (
	"github.com/JunLang-7/mall/adaptor"
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
}

func NewRouter(conf *config.Config, adaptor adaptor.IAdaptor, checkFunc func() error) *Router {
	return &Router{
		rootPath:  "/api/mall",
		conf:      conf,
		checkFunc: checkFunc,
	}
}

func (r *Router) Register(engine *gin.Engine) {
	engine.Use(AuthMiddleware(r.SpanFilter))
	SetupPprof(engine, "/debug/mall")
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
	root.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello world",
		})
	})
}
