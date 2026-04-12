package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) GetUserInfo(ctx *gin.Context) {
	// token auth
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := ctrl.user.GetUserInfo(ctx.Request.Context(), &common.AdminUser{})
	api.WriteResp(ctx, resp, errno)
}
