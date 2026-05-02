package goods

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/utils/tools"
	"github.com/gogf/gf/util/gconv"
	"gorm.io/gorm"
)

type ICourse interface {
	CreateCourse(ctx context.Context, req *do.CreateCourse) (int64, error)
	GetCourseInfo(ctx context.Context, id int64) (*model.CourseGood, error)
	UpdateCourse(ctx context.Context, req *do.UpdateCourse) error
	UpdateCourseStatus(ctx context.Context, req *do.UpdateCourseStatus) error
	ListCourse(ctx context.Context, req *do.CourseList) ([]*model.CourseGood, int64, error)
}

type Course struct {
	db *gorm.DB
}

func NewCourse(adaptor adaptor.IAdaptor) *Course {
	return &Course{
		db: adaptor.GetDB(),
	}
}

func (c *Course) CreateCourse(ctx context.Context, req *do.CreateCourse) (int64, error) {
	timeNow := time.Now()
	addCourse := &model.CourseGood{
		Name:           req.Name,
		CoverKey:       req.CoverKey,
		DetailCoverKey: req.DetailCoverKey,
		Detail:         req.Detail,
		CoursePrice:    req.CoursePrice,
		ServiceTime:    req.ServiceTime,
		LearnTime:      req.LearnTime,
		Status:         consts.IsDisable,
		Sort:           req.Sort,
		Features:       gconv.String(req.Features),
		UpdateStatus:   req.UpdateStatus,
		CreateAt:       timeNow,
		CreateBy:       req.UserID,
		UpdateAt:       timeNow,
		UpdateBy:       req.UserID,
	}
	err := c.db.WithContext(ctx).Create(addCourse).Error
	return addCourse.ID, err
}

func (c *Course) GetCourseInfo(ctx context.Context, id int64) (*model.CourseGood, error) {
	qs := query.Use(c.db).CourseGood
	return qs.WithContext(ctx).Where(qs.ID.Eq(id)).First()
}

func (c *Course) UpdateCourse(ctx context.Context, req *do.UpdateCourse) error {
	timeNow := time.Now()
	qs := query.Use(c.db).CourseGood
	updateMap := map[string]interface{}{
		qs.UpdateBy.ColumnName().String(): req.UserID,
		qs.UpdateAt.ColumnName().String(): timeNow,
	}
	if req.Name != "" {
		updateMap[qs.Name.ColumnName().String()] = req.Name
	}
	if req.CoursePrice != 0 {
		updateMap[qs.CoursePrice.ColumnName().String()] = req.CoursePrice
	}
	if req.ServiceTime != 0 {
		updateMap[qs.ServiceTime.ColumnName().String()] = req.ServiceTime
	}
	if req.LearnTime != 0 {
		updateMap[qs.LearnTime.ColumnName().String()] = req.LearnTime
	}
	if req.Sort != 0 {
		updateMap[qs.Sort.ColumnName().String()] = req.Sort
	}
	if req.UpdateStatus != 0 {
		updateMap[qs.UpdateStatus.ColumnName().String()] = req.UpdateStatus
	}
	if len(req.Features) != 0 {
		updateMap[qs.Features.ColumnName().String()] = req.Features
	}
	if req.CoverKey != "" {
		updateMap[qs.CoverKey.ColumnName().String()] = req.CoverKey
	}
	if req.DetailCoverKey != "" {
		updateMap[qs.DetailCoverKey.ColumnName().String()] = req.DetailCoverKey
	}
	if req.Detail != "" {
		updateMap[qs.Detail.ColumnName().String()] = req.Detail
	}
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(updateMap)
	return err
}

func (c *Course) UpdateCourseStatus(ctx context.Context, req *do.UpdateCourseStatus) error {
	qs := query.Use(c.db).CourseGood
	updateMap := map[string]interface{}{
		qs.UpdateBy.ColumnName().String(): req.UserID,
		qs.UpdateAt.ColumnName().String(): time.Now(),
		qs.Status.ColumnName().String():   req.Status,
	}
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(updateMap)
	return err
}

func (c *Course) ListCourse(ctx context.Context, req *do.CourseList) ([]*model.CourseGood, int64, error) {
	qs := query.Use(c.db).CourseGood
	q := qs.WithContext(ctx)
	if req.ID != 0 {
		q = q.Where(qs.ID.Eq(req.ID))
	}
	if req.NameKW != "" {
		q = q.Where(qs.Name.Like(tools.GetAllLike(req.NameKW)))
	}
	if req.CreateStartTime != 0 && req.CreateEndTime != 0 {
		q = q.Where(qs.CreateAt.Between(time.UnixMilli(req.CreateStartTime), time.UnixMilli(req.CreateEndTime)))
	}
	if req.UpdateStartTime != 0 && req.UpdateEndTime != 0 {
		q = q.Where(qs.UpdateAt.Between(time.UnixMilli(req.UpdateStartTime), time.UnixMilli(req.UpdateEndTime)))
	}
	if req.UpdateStatus != 0 {
		q = q.Where(qs.UpdateStatus.Eq(req.UpdateStatus))
	}
	if req.Status != 0 {
		q = q.Where(qs.Sort.Eq(req.Status))
	}
	return q.Order(qs.Status.Desc(), qs.ID.Desc()).FindByPage(req.GetOffset(), req.Limit)
}
