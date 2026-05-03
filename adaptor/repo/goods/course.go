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

	AddCatalog(ctx context.Context, req *do.AddCatalog) (int64, error)
	UpdateCatalog(ctx context.Context, req *do.UpdateCatalog) error
	DeleteCatalog(ctx context.Context, id int64) error
	UpdateCatalogSort(ctx context.Context, req *do.UpdateCatalogSort) error
	GetCatalogByCourseID(ctx context.Context, courseID int64) ([]*model.CourseCatalog, error)
	GetCourseLessons(ctx context.Context, courseID int64) ([]*model.CourseLesson, error)

	AddCatalogLesson(ctx context.Context, req *do.AddCatalogLesson) error
	RemoveCatalogLesson(ctx context.Context, req *do.RemoveCatalogLesson) error
	UpdateCatalogLesson(ctx context.Context, req *do.UpdateCatalogLesson) error
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

func (c *Course) AddCatalog(ctx context.Context, req *do.AddCatalog) (int64, error) {
	addCatalog := &model.CourseCatalog{
		CourseID: req.CourseID,
		Name:     req.Name,
		Level:    req.Level,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		UpdateAt: time.Now(),
		UpdateBy: req.UserID,
	}
	err := c.db.WithContext(ctx).Create(addCatalog).Error
	return addCatalog.ID, err

}

func (c *Course) UpdateCatalog(ctx context.Context, req *do.UpdateCatalog) error {
	qs := query.Use(c.db).CourseCatalog
	updateMap := map[string]interface{}{
		qs.UpdateBy.ColumnName().String(): req.UserID,
		qs.UpdateAt.ColumnName().String(): time.Now(),
		qs.Name.ColumnName().String():     req.Name,
	}
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(updateMap)
	return err
}

func (c *Course) DeleteCatalog(ctx context.Context, id int64) error {
	qs := query.Use(c.db).CourseCatalog
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(id)).Delete()
	return err
}

func (c *Course) UpdateCatalogSort(ctx context.Context, req *do.UpdateCatalogSort) error {
	timeNow := time.Now()
	cqs := query.Use(c.db).CourseCatalog
	lqs := query.Use(c.db).CourseLesson
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, catalog := range req.SortList {
			err := tx.Model(&model.CourseCatalog{}).Where(cqs.ID.Eq(catalog.ID)).Updates(map[string]interface{}{
				cqs.Sort.ColumnName().String():     catalog.Sort,
				cqs.ParentID.ColumnName().String(): catalog.ParentID,
				cqs.Level.ColumnName().String():    catalog.Level,
				cqs.UpdateAt.ColumnName().String(): timeNow,
				cqs.UpdateBy.ColumnName().String(): req.UserID,
			}).Error
			if err != nil {
				return err
			}
			for _, cLesson := range catalog.Lessons {
				err = tx.Model(&model.CourseLesson{}).Where(lqs.ID.Eq(cLesson.ID)).Updates(map[string]interface{}{
					lqs.Sort.ColumnName().String():     cLesson.Sort,
					lqs.UpdateBy.ColumnName().String(): req.UserID,
					lqs.UpdateAt.ColumnName().String(): timeNow,
				}).Error
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (c *Course) GetCourseLessons(ctx context.Context, courseID int64) ([]*model.CourseLesson, error) {
	qs := query.Use(c.db).CourseLesson
	return qs.WithContext(ctx).Where(qs.CourseGoodsID.Eq(courseID)).Order(qs.Sort).Find()
}

func (c *Course) GetCatalogByCourseID(ctx context.Context, courseID int64) ([]*model.CourseCatalog, error) {
	qs := query.Use(c.db).CourseCatalog
	return qs.WithContext(ctx).Where(qs.CourseID.Eq(courseID)).Order(qs.Sort).Find()
}

func (c *Course) AddCatalogLesson(ctx context.Context, req *do.AddCatalogLesson) error {
	timeNow := time.Now()
	qs := query.Use(c.db).CourseLesson
	addList := make([]*model.CourseLesson, 0, len(req.LessonIDs))
	for i, lessonID := range req.LessonIDs {
		addList = append(addList, &model.CourseLesson{
			CourseGoodsID: req.CourseID,
			CatalogID:     req.CatalogID,
			LessonID:      lessonID,
			Sort:          int32(i),
			ShowTime:      timeNow.Add(time.Minute * 5),
			UpdateAt:      timeNow,
			UpdateBy:      req.UserID,
		})
	}
	return qs.WithContext(ctx).CreateInBatches(addList, 100)
}

func (c *Course) RemoveCatalogLesson(ctx context.Context, req *do.RemoveCatalogLesson) error {
	qs := query.Use(c.db).CourseLesson
	_, err := qs.WithContext(ctx).Where(qs.ID.In(req.IDs...)).Delete()
	return err
}

func (c *Course) UpdateCatalogLesson(ctx context.Context, req *do.UpdateCatalogLesson) error {
	qs := query.Use(c.db).CourseLesson
	updateMap := map[string]interface{}{
		qs.UpdateBy.ColumnName().String(): req.UserID,
		qs.UpdateAt.ColumnName().String(): time.Now(),
	}
	if req.EnableTrial != 0 {
		updateMap[qs.EnableTrial.ColumnName().String()] = req.EnableTrial
	}
	if req.Name != "" {
		updateMap[qs.Name.ColumnName().String()] = req.Name
	}
	if req.ShowTime != 0 {
		updateMap[qs.ShowTime.ColumnName().String()] = time.UnixMilli(req.ShowTime)
	}
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(updateMap)
	return err
}
