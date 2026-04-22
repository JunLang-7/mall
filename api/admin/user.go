package admin

import (
	"context"
	"errors"

	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) GetUserInfo(ctx *gin.Context) {
	// token auth
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.GetUserInfoReq{}
	if err := ctx.ShouldBindQuery(req); err != nil || req.ID <= 0 {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg("invalid id"))
		return
	}
	resp, errno := ctrl.user.GetUserInfo(ctx.Request.Context(), user, req.ID)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) ListUsers(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.ListUsersReq{}
	resp, errno := ctrl.user.ListUsers(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) CreateUser(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.CreateUserReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg(err.Error()))
	}
	userId, errno := ctrl.user.CreateUser(ctx.Request.Context(), &common.AdminUser{}, req)
	api.WriteResp(ctx, map[string]int64{"id": userId}, errno)
}

func (ctrl *Ctrl) UpdateUser(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.UpdateUserReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg(err.Error()))
	}
	errno := ctrl.user.UpdateUser(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) UpdateUserStatus(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.UpdateUserStatusReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg(err.Error()))
	}
	errno := ctrl.user.UpdateUserStatus(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) DeleteUser(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.DeleteUserReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg(err.Error()))
		return
	}
	errno := ctrl.user.DeleteUser(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

// GetAdminUserByToken 获取管理员用户信息
func (ctrl *Ctrl) GetAdminUserByToken(ctx context.Context, token string) (*common.AdminUser, error) {
	adminUser, errno := ctrl.user.GetAdminUserByToken(ctx, token)
	if !errno.IsOK() {
		return nil, errors.New("failed to get admin user by token")
	}
	return adminUser, nil
}

// LarkBind 绑定飞书账号
func (ctrl *Ctrl) LarkBind(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.LarkBindReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg(err.Error()))
		return
	}
	errno := ctrl.user.LarkBind(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

// LarkUnbind 解绑飞书账号
func (ctrl *Ctrl) LarkUnbind(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	errno := ctrl.user.LarkUnbind(ctx.Request.Context(), user)
	api.WriteResp(ctx, nil, errno)
}

// AdminUserLogout 管理员用户登出
func (ctrl *Ctrl) AdminUserLogout(ctx *gin.Context) {
	adminUser := api.GetAdminUserFromCtx(ctx)
	errno := ctrl.user.AdminUserLogout(ctx.Request.Context(), adminUser)
	api.WriteResp(ctx, nil, errno)
}
