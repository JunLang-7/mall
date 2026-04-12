package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) GetUserInfo(ctx *gin.Context) {
	// token
	resp, errno := ctrl.user.GetUserInfo(ctx.Request.Context(), &common.AdminUser{})
	api.WriteResp(ctx, resp, errno)
}
