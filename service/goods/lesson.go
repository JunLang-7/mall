package goods

import (
	"context"

	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// CreateCategory 创建课程分类目录
func (s *Service) CreateCategory(ctx context.Context, req *dto.AddCategoryReq) (int64, common.Errno) {
	categoryID, err := s.lesson.CreateCategory(ctx, &do.AddCategory{
		Name:     req.Name,
		Level:    req.Level,
		ParentID: req.ParentID,
		Sort:     req.Sort,
	})
	if err != nil {
		logger.Error("CreateCategory CreateCategory error", zap.Error(err), zap.Any("req", req))
		return 0, *common.DataBaseErr.WithErr(err)
	}
	return categoryID, common.OK
}

// UpdateCategory 更新课程分类目录
func (s *Service) UpdateCategory(ctx context.Context, req *dto.UpdateCategoryReq) common.Errno {
	err := s.lesson.UpdateCategory(ctx, &do.UpdateCategory{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		logger.Error("UpdateCategory UpdateCategory error", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

// DeleteCategory 删除课程分类目录
func (s *Service) DeleteCategory(ctx context.Context, req *dto.DeleteCategoryReq) common.Errno {
	err := s.lesson.DeleteCategory(ctx, &do.DeleteCategory{
		ID: req.ID,
	})
	if err != nil {
		logger.Error("DeleteCategory DeleteCategory error", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

// ListCategory 获取课程分类目录列表
func (s *Service) ListCategory(ctx context.Context, req *dto.ListCategoryReq) ([]*dto.CategoryDto, common.Errno) {
	list, err := s.lesson.ListCategory(ctx, &do.ListCategory{
		Pager: req.Pager,
	})
	if err != nil {
		logger.Error("ListCategory ListCategory error", zap.Error(err), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	return lo.Map(list, func(item *model.LessonCategory, index int) *dto.CategoryDto {
		return &dto.CategoryDto{
			ID:       item.ID,
			Name:     item.Name,
			Sort:     item.Sort,
			ParentID: item.ParentID,
			Level:    item.Level,
		}
	}), common.OK
}

// CategorySort 更新课程分类目录排序
func (s *Service) CategorySort(ctx context.Context, sortList dto.UpdateCategorySortReq) common.Errno {
	updateSorts := make([]*do.UpdateSort, 0)
	lo.ForEach(sortList, func(item dto.UpdateSort, index int) {
		updateSorts = append(updateSorts, &do.UpdateSort{
			ID:       item.ID,
			Sort:     item.Sort,
			Level:    item.Level,
			ParentID: item.ParentID,
		})
	})
	err := s.lesson.CategorySort(ctx, updateSorts)
	if err != nil {
		logger.Error("CategorySort CategorySort error", zap.Error(err), zap.Any("sortList", sortList))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}
