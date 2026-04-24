package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

// GetTempSecret 获取对象存储临时密钥
func (ctrl *Ctrl) GetTempSecret(ctx *gin.Context) {
	req := &dto.GetTempSecretReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.storage.GetTempSecret(ctx, req)
	api.WriteResp(ctx, resp, errno)
}
