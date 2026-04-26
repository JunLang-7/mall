package goods

import (
	"context"
	"encoding/json"

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

// CreateLesson 创建课程
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

// UpdateLesson 更新课程
func (s *Service) UpdateLesson(ctx context.Context, user *common.AdminUser, req *dto.UpdateLessonReq) common.Errno {
	lesson, err := s.lesson.GetLessonByID(ctx, req.ID)
	if err != nil {
		logger.Error("UpdateLesson GetLessonByID error", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	if lesson == nil {
		return *common.ParamErr.WithMsg("invalid id")
	}
	err = s.lesson.UpdateLesson(ctx, &do.UpdateLesson{
		UserID:        user.UserID,
		ID:            req.ID,
		Name:          req.Name,
		Detail:        req.Detail,
		CategoryID:    req.CategoryID,
		VideoKey:      req.VideoKey,
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
		logger.Error("UpdateLesson UpdateLesson error", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	// 视频更换 删除cos中的旧视频文件
	if lesson.VideoKey != req.VideoKey {
		err := s.storage.DeleteFile(ctx, &do.DeleteFile{
			Keys: []string{lesson.VideoKey},
		})
		if err != nil {
			logger.Error("UpdateLesson DeleteFile error", zap.Error(err), zap.Any("req", req))
			return common.OK
		}
		err = s.upload.DeleteUploadFile(ctx, []string{lesson.VideoKey})
		if err != nil {
			logger.Error("UpdateLesson DeleteUploadFile error", zap.Error(err), zap.Any("req", req))
			return common.OK
		}
	}
	return common.OK
}

// UpdateLessonStatus 更新课程状态
func (s *Service) UpdateLessonStatus(ctx context.Context, user *common.AdminUser, req *dto.UpdateLessonStatusReq) common.Errno {
	err := s.lesson.UpdateLessonStatus(ctx, &do.UpdateLessonStatus{
		UserID: user.UserID,
		ID:     req.ID,
		Status: req.Status,
	})
	if err != nil {
		logger.Error("UpdateLessonStatus UpdateLessonStatus error", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

// MoveLesson 移动课程到其他分类
func (s *Service) MoveLesson(ctx context.Context, user *common.AdminUser, req *dto.MoveLessonReq) common.Errno {
	if len(req.LessonIDs) == 0 {
		errno := common.ParamErr
		errno.ErrMsg = "lesson_ids cannot be empty"
		return errno
	}
	err := s.lesson.MoveLesson(ctx, &do.MoveLesson{
		UserID:     user.UserID,
		LessonIDs:  req.LessonIDs,
		CategoryID: req.CategoryID,
	})
	if err != nil {
		logger.Error("MoveLesson MoveLesson error", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

// ListLesson 获取课程列表
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
	fileKeys := make([]string, 0)
	categoryNames := make(map[int64]string)
	userNames := make(map[int64]string)
	fileNames := make(map[string]string)

	lo.ForEach(list, func(item *model.Lesson, index int) {
		categoryIDs = append(categoryIDs, item.CategoryID)
		userIDs = append(userIDs, item.CreateBy, item.UpdateBy)
		fileKeys = append(fileKeys, item.VideoKey)
	})

	categoryIDs = lo.Uniq(categoryIDs)
	userIDs = lo.Uniq(userIDs)

	// 并发查询分类名称和用户名称
	pl := pool.NewPoolWithSize(3)
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
	pl.RunGo(func() {
		temp, err := s.storage.GetPreviewUrl(ctx, &do.GetPreviewUrl{
			Keys:   fileKeys,
			Expire: 6,
		})
		if err != nil {
			logger.Error("ListLesson GetFileNameMap error", zap.Error(err), zap.Any("req", req))
			return
		}
		fileNames = temp
	})
	pl.Wait()

	retList := make([]*dto.LessonDto, 0)
	lo.ForEach(list, func(item *model.Lesson, index int) {
		attachments := make([]dto.Attachment, 0)
		chapters := make([]dto.LessonChapter, 0)
		json.Unmarshal([]byte(item.Attachments), &attachments)
		json.Unmarshal([]byte(item.Chapters), &chapters)
		retList = append(retList, &dto.LessonDto{
			ID:            item.ID,
			Name:          item.Name,
			Detail:        item.Detail,
			Duration:      item.Duration,
			CategoryID:    item.CategoryID,
			CategoryName:  categoryNames[item.CategoryID],
			VideoKey:      item.VideoKey,
			VideoURL:      fileNames[item.VideoKey],
			VideoFileName: item.VideoFileName,
			Status:        item.Status,
			CreateBy:      item.CreateBy,
			UpdateBy:      item.UpdateBy,
			CreateAt:      item.CreateAt.UnixMilli(),
			UpdateAt:      item.UpdateAt.UnixMilli(),
			Attachments:   attachments,
			Chapters:      chapters,
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

func (s *Service) LessonInfo(ctx context.Context, req *dto.LessonInfoReq) (*dto.LessonDto, common.Errno) {
	resp, errno := s.ListLesson(ctx, &dto.ListLessonReq{
		ID: req.ID,
	})
	if !errno.IsOK() {
		logger.Error("LessonInfo ListLesson error", zap.Any("err", errno), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithMsg(errno.ErrMsg)
	}
	if resp == nil || len(resp.List) == 0 {
		logger.Error("LessonInfo ListLesson error", zap.Any("resp", resp))
		return nil, *common.ParamErr.WithMsg("invalid id")
	}
	return resp.List[0], common.OK
}
