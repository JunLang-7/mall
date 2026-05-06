package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) ListCustomerUsers(ctx *gin.Context) {
	req := &dto.AdminCustomerUserListReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.customer.AdminListCustomerUsers(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) UpdateCustomerUserStatus(ctx *gin.Context) {
	req := &dto.AdminCustomerUserStatusReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.customer.AdminUpdateCustomerStatus(ctx.Request.Context(), req)
	api.WriteResp(ctx, nil, errno)
}
