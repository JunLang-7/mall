package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) CreateCourse(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.CreateCourseReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	courseID, errno := ctrl.lesson.CreateCourse(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, map[string]int64{"id": courseID}, errno)
}

func (ctrl *Ctrl) GetCourseInfo(ctx *gin.Context) {
	req := &dto.CourseInfoReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.lesson.GetCourseInfo(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) UpdateCourse(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.UpdateCourseReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.UpdateCourse(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) UpdateCourseStatus(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.UpdateCourseStatusReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.UpdateCourseStatus(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) ListCourse(ctx *gin.Context) {
	req := &dto.CourseListReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, err := ctrl.lesson.ListCourse(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, err)
}
