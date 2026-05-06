package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) AdminOrderList(ctx *gin.Context) {
	req := &dto.OrderListReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.customer.AdminListOrders(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) AdminOrderInfo(ctx *gin.Context) {
	req := &dto.OrderInfoReq{}
	if err := ctx.ShouldBindQuery(req); err != nil || req.OrderID <= 0 {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg("invalid order_id"))
		return
	}
	resp, errno := ctrl.customer.AdminOrderInfo(ctx.Request.Context(), req.OrderID)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) AdminOrderStats(ctx *gin.Context) {
	req := &dto.AdminOrderStatsReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.customer.AdminOrderStats(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) AdminOrderRefund(ctx *gin.Context) {
	adminUser := api.GetAdminUserFromCtx(ctx)
	req := &dto.RefundOrderReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.customer.RefundOrder(ctx.Request.Context(), adminUser.UserID, req)
	api.WriteResp(ctx, nil, errno)
}
