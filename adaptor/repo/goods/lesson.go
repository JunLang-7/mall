package goods

import (
	"context"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/service/do"
	"gorm.io/gorm"
)

type ILesson interface {
	CreateCategory(ctx context.Context, req *do.AddCategory) (int64, error)
	UpdateCategory(ctx context.Context, req *do.UpdateCategory) error
	DeleteCategory(ctx context.Context, req *do.DeleteCategory) error
	ListCategory(ctx context.Context, req *do.ListCategory) ([]*model.LessonCategory, error)
	CategorySort(ctx context.Context, sortList []*do.UpdateSort) error
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
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Delete()
	return err
}

func (l *Lesson) ListCategory(ctx context.Context, req *do.ListCategory) ([]*model.LessonCategory, error) {
	qs := query.Use(l.db).LessonCategory
	return qs.WithContext(ctx).Order(qs.Sort).Find()
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
