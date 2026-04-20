package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/gin-gonic/gin"
)

// PermissionList 获取权限列表
func (ctrl *Ctrl) PermissionList(ctx *gin.Context) {
	// token auth
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := ctrl.perm.PermissionList(ctx.Request.Context())
	api.WriteResp(ctx, resp, errno)
}

// MyPermissionList 获取我的权限列表
func (ctrl *Ctrl) MyPermissionList(ctx *gin.Context) {
	// token auth
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := ctrl.perm.MyPermissionList(ctx.Request.Context(), user)
	api.WriteResp(ctx, resp, errno)
}
