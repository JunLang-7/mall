package customer

import (
	"context"
	"errors"

	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) MobileVerifyLogin(ctx *gin.Context) {
	req := &dto.MobileVerifyCodeLoginReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.MobileVerifyLogin(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) MobilePasswordLogin(ctx *gin.Context) {
	req := &dto.MobileLoginReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.MobilePasswordLogin(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) MobilePasswordReset(ctx *gin.Context) {
	req := &dto.MobilePasswordResetReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.user.MobilePasswordReset(ctx.Request.Context(), req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) GetUserInfo(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := ctrl.user.GetCustomerUserInfo(ctx.Request.Context(), user.UserID)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) ChangePassword(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.ChangePasswordReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.ChangePassword(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) GetCustomerUserByToken(ctx context.Context, token string) (*common.User, error) {
	user, errno := ctrl.user.GetCustomerUserByToken(ctx, token)
	if !errno.IsOK() {
		return nil, errors.New(errno.Msg)
	}
	return user, nil
}

func (ctrl *Ctrl) WechatQrCodeLogin(ctx *gin.Context) {
	api.WriteResp(ctx, &dto.WechatQrCodeResp{ExpireIn: 300, SceneToken: "mock", QrcodeURL: ""}, common.OK)
}

func (ctrl *Ctrl) WechatQrCodeStatus(ctx *gin.Context) {
	api.WriteResp(ctx, &dto.WechatQrCodeStatusResp{State: "pending", Purpose: "login", Message: "waiting"}, common.OK)
}

func (ctrl *Ctrl) WechatScanConfirm(ctx *gin.Context) {
	api.WriteResp(ctx, &dto.WechatQrCodeStatusResp{State: "confirmed", Purpose: "login", Message: "confirmed"}, common.OK)
}

func (ctrl *Ctrl) WechatQrCodeBind(ctx *gin.Context) {
	api.WriteResp(ctx, &dto.WechatQrCodeResp{ExpireIn: 300, SceneToken: "mock", QrcodeURL: ""}, common.OK)
}

func (ctrl *Ctrl) WechatUnbind(ctx *gin.Context) {
	api.WriteResp(ctx, nil, common.OK)
}
