package customer

import (
	"context"
	"errors"

	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) AppletLogin(ctx *gin.Context) {
	req := &dto.AppletLoginReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.AppletLogin(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

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

func (ctrl *Ctrl) GetChangePasswordSmsCode(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.ChangePasswordSmsCodeReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.user.SendChangePasswordSmsCode(ctx.Request.Context(), user.UserID, req.Ticket)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) GetCustomerUserByToken(ctx context.Context, token string) (*common.User, error) {
	user, errno := ctrl.user.GetCustomerUserByToken(ctx, token)
	if !errno.IsOK() {
		return nil, errors.New(errno.Msg)
	}
	return user, nil
}

func (ctrl *Ctrl) WechatQrCodeLogin(ctx *gin.Context) {
	req := &dto.WechatQrCodeReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.GetWechatQrCode(ctx.Request.Context(), req.Purpose)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) WechatQrCodeStatus(ctx *gin.Context) {
	sceneToken := ctx.Query("scene_token")
	if sceneToken == "" {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg("missing scene_token"))
		return
	}
	resp, errno := ctrl.user.GetWechatQrCodeStatus(ctx.Request.Context(), sceneToken)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) WechatScanConfirm(ctx *gin.Context) {
	req := &dto.WechatScanConfirmReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.ConfirmWechatScan(ctx.Request.Context(), req.SceneToken, req.Code)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) WechatQrCodeBind(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := ctrl.user.GetWechatQrCode(ctx.Request.Context(), "bind")
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) WechatUnbind(ctx *gin.Context) {
	api.WriteResp(ctx, nil, common.OK)
}
