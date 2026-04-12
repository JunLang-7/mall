package api

import (
	"net/http"

	"github.com/JunLang-7/mall/common"
	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	ErrMsg  string `json:"err_msg"`
	Data    any    `json:"data"`
}

func WriteResp(ctx *gin.Context, data any, errno common.Errno) {
	ctx.JSON(http.StatusOK, Resp{
		Code:    errno.Code,
		Message: errno.Msg,
		ErrMsg:  errno.ErrMsg,
		Data:    data,
	})
}
