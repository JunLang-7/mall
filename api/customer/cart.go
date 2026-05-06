package customer

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) AddCartGoods(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.AddGoodsReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	id, errno := ctrl.user.AddCartGoods(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, map[string]int64{"id": id}, errno)
}

func (ctrl *Ctrl) RemoveCartGoods(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.RemoveGoodsReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.user.RemoveCartGoods(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) ListCartGoods(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.ListCartGoodsReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.ListCartGoods(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, resp, errno)
}
