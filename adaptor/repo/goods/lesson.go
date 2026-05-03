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
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ILesson interface {
	CreateCategory(ctx context.Context, req *do.AddCategory) (int64, error)
	UpdateCategory(ctx context.Context, req *do.UpdateCategory) error
	DeleteCategory(ctx context.Context, req *do.DeleteCategory) error
	ListCategory(ctx context.Context, req *do.ListCategory) ([]*model.LessonCategory, error)
	GetChildCategoryIDs(ctx context.Context, parentIDs []int64) ([]int64, error)
	GetCategoryNameMap(ctx context.Context, ids []int64) (map[int64]string, error)
	CategorySort(ctx context.Context, sortList []*do.UpdateSort) error
	CreateLesson(ctx context.Context, req *do.CreateLesson) (int64, error)
	UpdateLesson(ctx context.Context, req *do.UpdateLesson) error
	UpdateLessonStatus(ctx context.Context, req *do.UpdateLessonStatus) error
	MoveLesson(ctx context.Context, req *do.MoveLesson) error
	ListLesson(ctx context.Context, req *do.ListLesson) ([]*model.Lesson, int64, error)
	GetLessonByID(ctx context.Context, id int64) (*model.Lesson, error)
	GetLessonByIDs(ctx context.Context, ids []int64) (map[int64]*model.Lesson, error)
}

type Lesson struct {
	db *gorm.DB
}

func NewLesson(adaptor adaptor.IAdaptor) *Lesson {
	return &Lesson{
		db: adaptor.GetDB(),
	}
}

func (l *Lesson) CreateCategory(ctx context.Context, req *do.AddCategory) (int64, error) {
	addObj := &model.LessonCategory{
		Name:     req.Name,
		ParentID: req.ParentID,
		Level:    req.Level,
		Sort:     req.Sort,
	}
	err := l.db.WithContext(ctx).Create(addObj).Error
	return addObj.ID, err
}

func (l *Lesson) UpdateCategory(ctx context.Context, req *do.UpdateCategory) error {
	qs := query.Use(l.db).LessonCategory
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Update(qs.Name, req.Name)
	return err
}

func (l *Lesson) DeleteCategory(ctx context.Context, req *do.DeleteCategory) error {
	qs := query.Use(l.db).LessonCategory
	_, err := qs.WithContext(ctx).Where(qs.ID.In(req.IDs...)).Delete()
	return err
}

func (l *Lesson) ListCategory(ctx context.Context, req *do.ListCategory) ([]*model.LessonCategory, error) {
	qs := query.Use(l.db).LessonCategory
	return qs.WithContext(ctx).Order(qs.Sort).Find()
}

func (l *Lesson) GetChildCategoryIDs(_ context.Context, parentIDs []int64) ([]int64, error) {
	var descendants []*model.LessonCategory
	rawSQL := `
	WITH RECURSIVE cte AS (
	  SELECT
		id,
		subject_id,
		parent_id,
		name,
		sort
	  FROM
		lesson_category
 	  WHERE
		id IN(?)
 	  UNION ALL
	  SELECT
		cd.id,
		cd.subject_id,
		cd.parent_id,
		cd.name,
		cd.sort
	  FROM
		lesson_category cd
	  INNER JOIN cte ON cd.parent_id = cte.id
	  )
	  SELECT * FROM cte
	`
	err := l.db.Exec(rawSQL, parentIDs).Scan(&descendants).Error
	if err != nil {
		return nil, err
	}
	categoryIDs := make([]int64, 0)
	lo.ForEach(descendants, func(item *model.LessonCategory, index int) {
		categoryIDs = append(categoryIDs, item.ID)
	})
	return categoryIDs, nil
}

func (l *Lesson) GetCategoryNameMap(ctx context.Context, ids []int64) (map[int64]string, error) {
	qs := query.Use(l.db).LessonCategory
	list, err := qs.WithContext(ctx).Select(qs.ID, qs.Name).Where(qs.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(list, func(item *model.LessonCategory) (int64, string) {
		return item.ID, item.Name
	}), nil
}

func (l *Lesson) CategorySort(ctx context.Context, sortList []*do.UpdateSort) error {
	qs := query.Use(l.db).LessonCategory
	return l.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range sortList {
			_, err := qs.WithContext(ctx).Where(qs.ID.Eq(v.ID)).Updates(map[string]interface{}{
				"sort":      v.Sort,
				"level":     v.Level,
				"parent_id": v.ParentID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (l *Lesson) CreateLesson(ctx context.Context, req *do.CreateLesson) (int64, error) {
	timeNow := time.Now()
	addObj := &model.Lesson{
		Name:          req.Name,
		Detail:        req.Detail,
		CategoryID:    req.CategoryID,
		VideoKey:      req.VideoKey,
		VideoFileName: req.VideoFileName,
		Duration:      req.Duration,
		Attachments:   gconv.String(req.Attachments),
		Chapters:      gconv.String(req.Chapters),
		Status:        consts.IsEnable,
		CreateBy:      req.UserID,
		UpdateBy:      req.UserID,
		CreateAt:      timeNow,
		UpdateAt:      timeNow,
	}
	err := l.db.WithContext(ctx).Create(addObj).Error
	return addObj.ID, err
}

func (l *Lesson) UpdateLesson(ctx context.Context, req *do.UpdateLesson) error {
	qs := query.Use(l.db).Lesson
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).UpdateSimple(
		qs.Name.Value(req.Name),
		qs.Detail.Value(req.Detail),
		qs.CategoryID.Value(req.CategoryID),
		qs.VideoKey.Value(req.VideoKey),
		qs.VideoFileName.Value(req.VideoFileName),
		qs.Duration.Value(req.Duration),
		qs.Attachments.Value(gconv.String(req.Attachments)),
		qs.Chapters.Value(gconv.String(req.Chapters)),
		qs.UpdateBy.Value(req.UserID),
		qs.UpdateAt.Value(time.Now()),
	)
	return err
}

func (l *Lesson) UpdateLessonStatus(ctx context.Context, req *do.UpdateLessonStatus) error {
	qs := query.Use(l.db).Lesson
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).UpdateSimple(
		qs.Status.Value(req.Status),
		qs.UpdateAt.Value(time.Now()),
		qs.UpdateBy.Value(req.UserID),
	)
	return err
}

func (l *Lesson) MoveLesson(ctx context.Context, req *do.MoveLesson) error {
	qs := query.Use(l.db).Lesson
	_, err := qs.WithContext(ctx).Where(qs.ID.In(req.LessonIDs...)).UpdateSimple(
		qs.CategoryID.Value(req.CategoryID),
		qs.UpdateAt.Value(time.Now()),
		qs.UpdateBy.Value(req.UserID),
	)
	return err
}

func (l *Lesson) ListLesson(ctx context.Context, req *do.ListLesson) ([]*model.Lesson, int64, error) {
	qs := query.Use(l.db).Lesson
	tx := qs.WithContext(ctx)
	if req.ID > 0 {
		tx = tx.Where(qs.ID.Eq(req.ID))
	}
	if len(req.ExcludeCourseIDs) > 0 {
		tx = tx.Where(qs.ID.NotIn(req.ExcludeCourseIDs...))
	}
	if req.NameKw != "" {
		tx = tx.Where(qs.Name.Like(tools.GetAllLike(req.NameKw)))
	}
	if len(req.CategoryIDs) > 0 {
		tx = tx.Where(qs.CategoryID.In(req.CategoryIDs...))
	}
	if req.Status != 0 {
		tx = tx.Where(qs.Status.Eq(req.Status))
	}
	if req.StartCreateTime > 0 && req.EndCreateTime > 0 {
		tx = tx.Where(qs.CreateAt.Between(time.UnixMilli(req.StartCreateTime), time.UnixMilli(req.EndCreateTime)))
	}
	if req.StartUpdateTime > 0 && req.EndUpdateTime > 0 {
		tx = tx.Where(qs.UpdateAt.Between(time.UnixMilli(req.EndUpdateTime), time.UnixMilli(req.EndUpdateTime)))
	}
	return tx.Order(qs.Status.Desc(), qs.ID.Desc()).FindByPage(req.GetOffset(), req.Limit)
}

func (l *Lesson) GetLessonByID(ctx context.Context, id int64) (*model.Lesson, error) {
	qs := query.Use(l.db).Lesson
	return qs.WithContext(ctx).Where(qs.ID.Eq(id)).First()
}

func (l *Lesson) GetLessonByIDs(ctx context.Context, ids []int64) (map[int64]*model.Lesson, error) {
	qs := query.Use(l.db).Lesson
	list, err := qs.WithContext(ctx).Where(qs.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(list, func(item *model.Lesson) (int64, *model.Lesson) {
		return item.ID, item
	}), nil
}
