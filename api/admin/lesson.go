package admin

import (
	"github.com/JunLang-7/mall/api"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/gin-gonic/gin"
)

// CreateCategory 创建课程分类目录
func (ctrl *Ctrl) CreateCategory(ctx *gin.Context) {
	req := &dto.AddCategoryReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.lesson.CreateCategory(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

// UpdateCategory 更新课程分类目录
func (ctrl *Ctrl) UpdateCategory(ctx *gin.Context) {
	req := &dto.UpdateCategoryReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.UpdateCategory(ctx.Request.Context(), req)
	api.WriteResp(ctx, nil, errno)
}

// DeleteCategory 删除课程分类目录
func (ctrl *Ctrl) DeleteCategory(ctx *gin.Context) {
	req := &dto.DeleteCategoryReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.DeleteCategory(ctx.Request.Context(), req)
	api.WriteResp(ctx, nil, errno)
}

// ListCategory 获取课程分类目录列表
func (ctrl *Ctrl) ListCategory(ctx *gin.Context) {
	req := &dto.ListCategoryReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.lesson.ListCategory(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

// CategorySorts 更新课程分类目录排序
func (ctrl *Ctrl) CategorySorts(ctx *gin.Context) {
	req := dto.UpdateCategorySortReq{}
	if err := ctx.BindJSON(&req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.CategorySort(ctx.Request.Context(), req)
	api.WriteResp(ctx, nil, errno)
}

// CreateLesson 创建课程
func (ctrl *Ctrl) CreateLesson(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.CreateLessonReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	id, errno := ctrl.lesson.CreateLesson(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, id, errno)
}

// UpdateLesson 更新课程
func (ctrl *Ctrl) UpdateLesson(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.UpdateLessonReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.UpdateLesson(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

// UpdateLessonStatus 更新课程状态
func (ctrl *Ctrl) UpdateLessonStatus(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.UpdateLessonStatusReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.UpdateLessonStatus(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

// MoveLesson 移动课程到其他分类
func (ctrl *Ctrl) MoveLesson(ctx *gin.Context) {
	user := api.GetAdminUserFromCtx(ctx)
	req := &dto.MoveLessonReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	errno := ctrl.lesson.MoveLesson(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

// ListLesson 获取课程列表
func (ctrl *Ctrl) ListLesson(ctx *gin.Context) {
	req := &dto.ListLessonReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.lesson.ListLesson(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (ctrl *Ctrl) LessonInfo(ctx *gin.Context) {
	req := &dto.LessonInfoReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, *common.ParamErr.WithErr(err))
		return
	}
	resp, errno := ctrl.lesson.LessonInfo(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}
