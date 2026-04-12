package admin

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/gin-gonic/gin"
)

type Ctrl struct {
	adaptor adaptor.IAdaptor
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{adaptor: adaptor}
}

func (ctrl *Ctrl) Hello(ctx *gin.Context) {
	api.WriteResp(ctx, "Hello World", common.OK)
}
