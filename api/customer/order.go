package customer

import (
	"net/http"
	"strconv"

	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) CalcOrderFee(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.OrderCalcFeeReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.CalcOrderFee(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) PayNow(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.OrderPayNowReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.PayNow(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) PayLater(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.OrderPayLaterReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.PayLater(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) CancelOrder(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.CancelOrderReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.user.CancelOrder(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) OrderList(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.OrderListReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.ListOrders(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) OrderInfo(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.OrderInfoReq{}
	if err := ctx.ShouldBindQuery(req); err != nil || req.OrderID <= 0 {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg("invalid order_id"))
		return
	}
	resp, errno := ctrl.user.GetOrderInfo(ctx.Request.Context(), user.UserID, req.OrderID)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) WechatPaymentCallback(ctx *gin.Context) {
	orderID, _ := strconv.ParseInt(ctx.Query("order_id"), 10, 64)
	errno := ctrl.user.WechatNotifySuccess(ctx.Request.Context(), orderID)
	if !errno.IsOK() {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": "FAIL", "message": errno.Msg})
		return
	}
	ctx.JSON(http.StatusOK, ctrl.user.WechatNotifyResponse())
}

func (ctrl *Ctrl) WechatRefundCallback(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ctrl.user.WechatNotifyResponse())
}
