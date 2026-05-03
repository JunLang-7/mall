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

func (ctrl *Ctrl) AddCatalog(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.AddCatalogReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	catalogID, errno := ctrl.lesson.AddCatalog(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, map[string]int64{"id": catalogID}, errno)
}

func (ctrl *Ctrl) UpdateCatalog(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.UpdateCatalogReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.UpdateCatalog(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) DeleteCatalog(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.DeleteCatalogReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.DeleteCatalog(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) UpdateCatalogSort(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	sorts := make([]dto.UpdateCatalogSortDto, 0)
	if err := ctx.BindJSON(&sorts); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.UpdateCatalogSort(ctx.Request.Context(), user, sorts)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) GetCatalogInfo(ctx *gin.Context) {
	req := &dto.CatalogInfoReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.lesson.GetCatalogInfo(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) AddCatalogLesson(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.AddCatalogLessonReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.AddCatalogLesson(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) RemoveCatalogLesson(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.RemoveCatalogLessonReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
	}
	errno := ctrl.lesson.RemoveCatalogLesson(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (ctrl *Ctrl) UpdateCatalogLesson(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.UpdateCatalogLessonReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.UpdateCatalogLesson(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}
