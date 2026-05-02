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

func (s *Service) CreateCourse(ctx context.Context, user *common.AdminUser, req *dto.CreateCourseReq) (int64, common.Errno) {
	courseID, err := s.course.CreateCourse(ctx, &do.CreateCourse{
		UserID:         user.UserID,
		Name:           req.Name,
		CoursePrice:    req.CoursePrice,
		ServiceTime:    req.ServiceTime,
		LearnTime:      req.LearnTime,
		Sort:           req.Sort,
		Features:       req.Features,
		UpdateStatus:   req.UpdateStatus,
		CoverKey:       req.CoverKey,
		DetailCoverKey: req.DetailCoverKey,
		Detail:         req.Detail,
	})
	if err != nil {
		logger.Error("CreateCourse CreateCourse error", zap.Error(err), zap.Any("req", req))
		return 0, *common.DataBaseErr.WithErr(err)
	}
	return courseID, common.OK
}

func (s *Service) GetCourseInfo(ctx context.Context, req *dto.CourseInfoReq) (*dto.CourseDto, common.Errno) {
	course, err := s.course.GetCourseInfo(ctx, req.ID)
	if err != nil {
		logger.Error("GetCourseInfo error", zap.Error(err), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}
	var feature []string
	_ = json.Unmarshal([]byte(course.Features), &feature)

	fileUrlMap, err := s.storage.GetPreviewUrl(ctx, &do.GetPreviewUrl{
		Keys:   []string{course.CoverKey, course.DetailCoverKey},
		Expire: 6,
	})
	if err != nil {
		logger.Error("GetCourseInfo error", zap.Error(err), zap.Any("req", req))
		return nil, *common.ServerErr.WithErr(err)
	}

	return &dto.CourseDto{
		CreateUpdateName: common.CreateUpdateName{
			CreateName: course.Name,
			UpdateName: course.Name,
		},
		ID:             course.ID,
		Name:           course.Name,
		CoursePrice:    course.CoursePrice,
		ServiceTime:    course.ServiceTime,
		LearnTime:      course.LearnTime,
		Sort:           course.Sort,
		Status:         course.Status,
		Features:       feature,
		UpdateStatus:   course.UpdateStatus,
		CoverKey:       course.CoverKey,
		CoverURL:       fileUrlMap[course.CoverKey],
		DetailCoverKey: fileUrlMap[course.DetailCoverKey],
		Detail:         course.Detail,
		CreateBy:       course.CreateBy,
		CreateAt:       course.CreateAt.UnixMilli(),
		UpdateBy:       course.UpdateBy,
		UpdateAt:       course.UpdateAt.UnixMilli(),
	}, common.OK
}

func (s *Service) UpdateCourse(ctx context.Context, user *common.AdminUser, req *dto.UpdateCourseReq) common.Errno {
	err := s.course.UpdateCourse(ctx, &do.UpdateCourse{
		UserID:         user.UserID,
		ID:             req.ID,
		Name:           req.Name,
		CoursePrice:    req.CoursePrice,
		ServiceTime:    req.ServiceTime,
		LearnTime:      req.LearnTime,
		Sort:           req.Sort,
		Features:       req.Features,
		UpdateStatus:   req.UpdateStatus,
		CoverKey:       req.CoverKey,
		DetailCoverKey: req.DetailCoverKey,
		Detail:         req.Detail,
	})
	if err != nil {
		logger.Error("UpdateCourse UpdateCourse error", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) UpdateCourseStatus(ctx context.Context, user *common.AdminUser, req *dto.UpdateCourseStatusReq) common.Errno {
	err := s.course.UpdateCourseStatus(ctx, &do.UpdateCourseStatus{
		UserID: user.UserID,
		ID:     req.ID,
		Status: req.Status,
	})
	if err != nil {
		logger.Error("UpdateCourseStatus UpdateCourseStatus", zap.Error(err), zap.Any("req", req))
		return *common.DataBaseErr.WithErr(err)
	}
	return common.OK
}

func (s *Service) ListCourse(ctx context.Context, req *dto.CourseListReq) (*dto.CourseListResp, common.Errno) {
	list, total, err := s.course.ListCourse(ctx, &do.CourseList{
		Pager:           req.Pager,
		ID:              req.ID,
		NameKW:          req.NameKW,
		CreateStartTime: req.CreateStartTime,
		CreateEndTime:   req.CreateEndTime,
		UpdateStartTime: req.UpdateStartTime,
		UpdateEndTime:   req.UpdateEndTime,
		UpdateStatus:    req.UpdateStatus,
		Status:          req.Status,
	})
	if err != nil {
		logger.Error("ListCourse ListCourse error", zap.Error(err), zap.Any("req", req))
		return nil, *common.DataBaseErr.WithErr(err)
	}

	fileKeys := make([]string, 0)
	fileUrlMap := make(map[string]string)
	userNameMap := make(map[int64]string)
	userIDs := make([]int64, 0)
	lo.ForEach(list, func(item *model.CourseGood, index int) {
		fileKeys = append(fileKeys, item.CoverKey, item.DetailCoverKey)
		userIDs = append(userIDs, item.CreateBy, item.UpdateBy)
	})

	pl := pool.NewPoolWithSize(2)
	defer pl.Release()

	pl.RunGo(func() {
		tempMap, err := s.storage.GetPreviewUrl(ctx, &do.GetPreviewUrl{
			Keys:   fileKeys,
			Expire: 6,
		})
		if err != nil {
			logger.Error("ListCourse GetPreviewUrl error", zap.Error(err), zap.Any("req", req))
			return
		}
		fileUrlMap = tempMap
	})

	pl.RunGo(func() {
		tempMap, err := s.user.GetUserNameMap(ctx, userIDs)
		if err != nil {
			logger.Error("ListCourse GetUserNameMap error", zap.Error(err), zap.Any("req", req))
			return
		}
		userNameMap = tempMap
	})

	pl.Wait()

	retList := make([]*dto.CourseDto, 0, len(list))
	for _, course := range list {
		features := make([]string, 0)
		_ = json.Unmarshal([]byte(course.Features), &features)
		retList = append(retList, &dto.CourseDto{
			ID:             course.ID,
			Name:           course.Name,
			CoursePrice:    course.CoursePrice,
			ServiceTime:    course.ServiceTime,
			LearnTime:      course.LearnTime,
			Sort:           course.Sort,
			Status:         course.Status,
			Features:       features,
			UpdateStatus:   course.UpdateStatus,
			CoverKey:       course.CoverKey,
			CoverURL:       fileUrlMap[course.CoverKey],
			DetailCoverKey: course.DetailCoverKey,
			DetailCOverURL: fileUrlMap[course.DetailCoverKey],
			Detail:         course.Detail,
			CreateBy:       course.CreateBy,
			CreateAt:       course.CreateAt.UnixMilli(),
			UpdateBy:       course.UpdateBy,
			UpdateAt:       course.UpdateAt.UnixMilli(),
			CreateUpdateName: common.CreateUpdateName{
				CreateName: userNameMap[course.CreateBy],
				UpdateName: userNameMap[course.UpdateBy],
			},
		})
	}
	return &dto.CourseListResp{
		Pager: req.Pager,
		List:  retList,
		Total: total,
	}, common.OK
}
