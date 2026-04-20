package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
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

// CreatePermission 创建权限
func (ctrl *Ctrl) CreatePermission(ctx *gin.Context) {
	// token auth
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.PermissionCreateReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	permID, errno := ctrl.perm.CreatePermission(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, map[string]interface{}{
		"id": permID,
	}, errno)
}

// UpdatePermission 批量更新权限
func (ctrl *Ctrl) UpdatePermission(ctx *gin.Context) {
	// token auth
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.PermissionUpdateResp{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.perm.UpdatePermissions(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

// DeletePermission 删除权限
func (ctrl *Ctrl) DeletePermission(ctx *gin.Context) {
	// token auth
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.PermissionDeleteReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.perm.DeletePermission(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}
