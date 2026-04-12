package admin

import (
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
