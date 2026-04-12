package admin

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/admin"
	"github.com/gin-gonic/gin"
)

type Ctrl struct {
	adaptor adaptor.IAdaptor
	hello   *admin.Service
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		hello:   admin.NewService(adaptor),
	}
}

func (ctrl *Ctrl) Hello(ctx *gin.Context) {
	resp, errno := ctrl.hello.Hello(ctx.Request.Context(), &common.AdminUser{}, nil)
	api.WriteResp(ctx, resp, errno)
}
