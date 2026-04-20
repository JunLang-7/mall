package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

// AddRole 创建角色
func (ctrl *Ctrl) AddRole(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.AddRoleReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	id, errno := ctrl.role.CreateRole(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, map[string]interface{}{"id": id}, errno)
}

// UpdateRole 更新角色
func (ctrl *Ctrl) UpdateRole(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.UpdateRoleReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.role.UpdateRole(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

// ListRole 角色列表
func (ctrl *Ctrl) ListRole(ctx *gin.Context) {
	req := &dto.ListRoleReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.role.ListRole(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

// MyRoles 我的角色列表	
func (ctrl *Ctrl) MyRoles(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	roles, errno := ctrl.role.GetMyRoles(ctx.Request.Context(), user)
	api.WriteResp(ctx, roles, errno)
}

// SetRolePerms 设置角色权限
func (ctrl *Ctrl) SetRolePerms(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.SetRolePermReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.role.SetRolePerms(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}
