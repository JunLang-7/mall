package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

// GetSmsCodeCaptcha 获取短信验证码的滑动验证码
func (ctrl *Ctrl) GetSmsCodeCaptcha(ctx *gin.Context) {
	req := &dto.GetVerifyCaptchaReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.GetSlideCaptcha(ctx.Request.Context())
	api.WriteResp(ctx, resp, errno)
}

// CheckSmsCodeCaptcha 校验短信验证码的滑动验证码
func (ctrl *Ctrl) CheckSmsCodeCaptcha(ctx *gin.Context) {
	req := &dto.CheckCaptchaReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.CheckSlideCaptcha(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

// GetSmsCodeVerify 获取短信验证码
func (ctrl *Ctrl) GetSmsCodeVerify(ctx *gin.Context) {
	req := &dto.GetSmsCodeVerifyReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.user.GetSmsCodeVerify(ctx.Request.Context(), req)
	api.WriteResp(ctx, nil, errno)
}

// MobilePasswordLogin 手机号密码登录
func (ctrl *Ctrl) MobilePasswordLogin(ctx *gin.Context) {
	req := &dto.MobileLoginReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.MobilePasswordLogin(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) MobileVerifyLogin(ctx *gin.Context) {}

// LarkQrCodeLogin 飞书扫码登录
func (ctrl *Ctrl) LarkQrCodeLogin(ctx *gin.Context) {
	req := &dto.LarkQrCodeLoginReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.LarkQrCodeLogin(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}
