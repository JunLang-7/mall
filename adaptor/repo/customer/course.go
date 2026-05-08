package customer

import (
	"context"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"github.com/JunLang-7/mall/utils/tools"
	"gorm.io/gorm"
)

type ICourse interface {
	ListCourse(ctx context.Context, req *do.CourseListReq) ([]*model.CourseGood, int64, error)
	GetCourseInfo(ctx context.Context, req *do.CourseInfoReq) (*model.CourseGood, error)
	GetCourseByID(ctx context.Context, id int64) (*model.CourseGood, error)
	GetCatalogs(ctx context.Context, courseID int64) ([]*model.CourseCatalog, error)
	GetCourseLessons(ctx context.Context, courseID int64) ([]*model.CourseLesson, error)
}

type CourseRepo struct {
	db *gorm.DB
}

func NewCourseRepo(adaptor adaptor.IAdaptor) *CourseRepo {
	return &CourseRepo{db: adaptor.GetDB()}
}

func (r *CourseRepo) ListCourse(ctx context.Context, req *do.CourseListReq) ([]*model.CourseGood, int64, error) {
	qs := query.Use(r.db).CourseGood
	tx := qs.WithContext(ctx).Where(qs.Status.Eq(consts.IsEnable))
	if req.ID > 0 {
		tx = tx.Where(qs.ID.Eq(req.ID))
	}
	if req.NameKW != "" {
		tx = tx.Where(qs.Name.Like(tools.GetAllLike(req.NameKW)))
	}
	return tx.Order(qs.Sort.Asc(), qs.ID.Desc()).FindByPage(req.GetOffset(), req.Limit)
}

func (r *CourseRepo) GetCourseInfo(ctx context.Context, req *do.CourseInfoReq) (*model.CourseGood, error) {
	qs := query.Use(r.db).CourseGood
	tx := qs.WithContext(ctx).Where(qs.ID.Eq(req.ID))
	if req.Status != 0 {
		tx = tx.Where(qs.Status.Eq(req.Status))
	}
	return tx.First()
}

func (r *CourseRepo) GetCourseByID(ctx context.Context, id int64) (*model.CourseGood, error) {
	qs := query.Use(r.db).CourseGood
	return qs.WithContext(ctx).Where(qs.ID.Eq(id)).First()
}

func (r *CourseRepo) GetCatalogs(ctx context.Context, courseID int64) ([]*model.CourseCatalog, error) {
	qs := query.Use(r.db).CourseCatalog
	return qs.WithContext(ctx).Where(qs.CourseID.Eq(courseID)).Order(qs.Sort.Asc()).Find()
}

func (r *CourseRepo) GetCourseLessons(ctx context.Context, courseID int64) ([]*model.CourseLesson, error) {
	qs := query.Use(r.db).CourseLesson
	return qs.WithContext(ctx).Where(qs.CourseGoodsID.Eq(courseID)).Order(qs.Sort.Asc()).Find()
}
