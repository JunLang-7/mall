package goods

import (
	"context"

	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/common"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/service/dto"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/JunLang-7/mall/utils/pool"
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
	if len(req.IDs) == 0 {
		errno := common.ParamErr
		errno.ErrMsg = "ids cannot be empty"
		return errno
	}

	err := s.lesson.DeleteCategory(ctx, &do.DeleteCategory{
		IDs: req.IDs,
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

func (s *Service) CreateLesson(ctx context.Context, user *common.AdminUser, req *dto.CreateLessonReq) (int64, common.Errno) {
	lessonID, err := s.lesson.CreateLesson(ctx, &do.CreateLesson{
		UserID:        user.UserID,
		Name:          req.Name,
		Detail:        req.Detail,
		CategoryID:    req.CategoryID,
		VideoKey:      req.VideoKey, // TODO: 先查询回来，比较key是否一致，若不一致，进行删除cos文件
		VideoFileName: req.VideoFileName,
		Duration:      req.Duration,
		Attachments: lo.Map(req.Attachments, func(item dto.Attachment, index int) do.Attachment {
			return do.Attachment{
				FileKey:    item.FileKey,
				OriginName: item.OriginName,
			}
		}),
		Chapters: lo.Map(req.Chapters, func(item dto.LessonChapter, index int) do.LessonChapter {
			return do.LessonChapter{
				Name:          item.Name,
				BeginPosition: item.BeginPosition,
				EndPosition:   item.EndPosition,
			}
		}),
	})
	if err != nil {
		logger.Error("CreateLesson CreateLesson error", zap.Error(err), zap.Any("req", req))
		return 0, *common.DataBaseErr.WithErr(err)
	}
	return lessonID, common.OK
}

func (s *Service) ListLesson(ctx context.Context, req *dto.ListLessonReq) (*dto.ListLessonResp, common.Errno) {
	categoryIDs := make([]int64, 0)
	if req.CategoryID != 0 {
		categoryIDs = append(categoryIDs, req.CategoryID)
		if req.OnView {
			tempIDs, err := s.lesson.GetChildCategoryIDs(ctx, []int64{req.CategoryID})
			if err != nil {
				logger.Error("ListLesson GetChildCategoryIDs error", zap.Error(err), zap.Any("req", req))
				return nil, *common.DataBaseErr.WithErr(err)
			}
			categoryIDs = append(categoryIDs, tempIDs...)
		}
	}
	list, total, err := s.lesson.ListLesson(ctx, &do.ListLesson{
		Pager:           req.Pager,
		ID:              req.ID,
		NameKw:          req.NameKw,
		CategoryIDs:     categoryIDs,
		Status:          req.Status,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		StartUpdateTime: req.BeginUpdateTime,
		EndUpdateTime:   req.EndUpdateTime,
	})
	if err != nil {
		logger.Error("ListLesson ListLesson error", zap.Error(err), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	categoryIDs = make([]int64, 0)
	userIDs := make([]int64, 0)
	categoryNames := make(map[int64]string)
	userNames := make(map[int64]string)

	lo.ForEach(list, func(item *model.Lesson, index int) {
		categoryIDs = append(categoryIDs, item.CategoryID)
		userIDs = append(userIDs, item.CreateBy, item.UpdateBy)
	})

	categoryIDs = lo.Uniq(categoryIDs)
	userIDs = lo.Uniq(userIDs)

	pl := pool.NewPoolWithSize(2)
	defer pl.Release()
	pl.RunGo(func() {
		temp, err := s.lesson.GetCategoryNameMap(ctx, categoryIDs)
		if err != nil {
			logger.Error("ListLesson GetCategoryNameMap error", zap.Error(err), zap.Any("req", req))
			return
		}
		categoryNames = temp
	})
	pl.RunGo(func() {
		temp, err := s.user.GetUserNameMap(ctx, userIDs)
		if err != nil {
			logger.Error("ListLesson GetUserNameMap error", zap.Error(err), zap.Any("req", req))
			return
		}
		userNames = temp
	})
	pl.Wait()

	retList := make([]*dto.LessonDto, 0)
	lo.ForEach(list, func(item *model.Lesson, index int) {
		retList = append(retList, &dto.LessonDto{
			ID:           item.ID,
			Name:         item.Name,
			Detail:       item.Detail,
			Duration:     item.Duration,
			CategoryID:   item.CategoryID,
			CategoryName: categoryNames[item.CategoryID],
			Status:       item.Status,
			CreateBy:     item.CreateBy,
			UpdateBy:     item.UpdateBy,
			CreateAt:     item.CreateAt.UnixMilli(),
			UpdateAt:     item.UpdateAt.UnixMilli(),
			CreateUpdateName: common.CreateUpdateName{
				CreateName: userNames[item.CreateBy],
				UpdateName: userNames[item.UpdateBy],
			},
		})
	})
	return &dto.ListLessonResp{
		List:  retList,
		Total: total,
		Pager: req.Pager,
	}, common.OK
}
