package customer

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

func (ctrl *Ctrl) ListCourse(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.CourseListReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	userID := int64(0)
	if user != nil {
		userID = user.UserID
	}
	resp, errno := ctrl.user.ListCustomerCourse(ctx.Request.Context(), userID, req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) CourseDetail(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.CourseInfoReq{}
	if err := ctx.ShouldBindQuery(req); err != nil || req.ID <= 0 {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg("invalid id"))
		return
	}
	userID := int64(0)
	if user != nil {
		userID = user.UserID
	}
	resp, errno := ctrl.user.GetCustomerCourseDetail(ctx.Request.Context(), userID, req.ID)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) LessonInfo(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.LessonInfoReq{}
	if err := ctx.ShouldBindQuery(req); err != nil || req.ID <= 0 {
		api.WriteResp(ctx, nil, *common.ParamErr.WithMsg("invalid lesson_id"))
		return
	}
	resp, errno := ctrl.user.GetLessonInfo(ctx.Request.Context(), user.UserID, req.ID)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) LessonLearnInfo(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.LessonLearnInfoReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.GetLessonLearnInfo(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) LessonLearnReport(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &dto.LessonLearnReportReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.user.ReportLessonLearn(ctx.Request.Context(), user.UserID, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) PurchasedCourseList(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &common.Pager{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.ListPurchasedCourse(ctx.Request.Context(), user.UserID, *req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) ContinueLearnList(ctx *gin.Context) {
	user := api.GetUserFromCtx(ctx)
	req := &common.Pager{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.user.ListContinueLearn(ctx.Request.Context(), user.UserID, *req)
	api.WriteResp(ctx, resp, errno)
}
