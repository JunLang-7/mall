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
