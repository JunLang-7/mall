package customer

import (
	"context"
	"time"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/model"
	"github.com/JunLang-7/mall/adaptor/repo/query"
	"github.com/JunLang-7/mall/consts"
	"github.com/JunLang-7/mall/service/do"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ILesson interface {
	GetLessonInfo(ctx context.Context, lessonID int64) (*model.Lesson, error)
	GetCourseLesson(ctx context.Context, lessonID int64) (*model.CourseLesson, error)
	GetLessonLearnProgress(ctx context.Context, req *do.LessonLearnProgressReq) (*model.LessonLearnProgress, error)
	GetUserLearnProgresses(ctx context.Context, userID int64, offset, limit int) ([]*model.LessonLearnProgress, int64, error)
	UpsertLessonLearnProgress(ctx context.Context, req *do.LessonLearnProgressUpdate) error
	CreateLessonLearnRecord(ctx context.Context, req *do.LessonLearnRecordCreate) error
}

type LessonRepo struct {
	db *gorm.DB
}

func NewLessonRepo(adaptor adaptor.IAdaptor) *LessonRepo {
	return &LessonRepo{db: adaptor.GetDB()}
}

func (r *LessonRepo) GetLessonInfo(ctx context.Context, lessonID int64) (*model.Lesson, error) {
	qs := query.Use(r.db).Lesson
	return qs.WithContext(ctx).Where(qs.ID.Eq(lessonID)).First()
}

func (r *LessonRepo) GetCourseLesson(ctx context.Context, lessonID int64) (*model.CourseLesson, error) {
	qs := query.Use(r.db).CourseLesson
	return qs.WithContext(ctx).Where(qs.LessonID.Eq(lessonID)).First()
}

func (r *LessonRepo) GetLessonLearnProgress(ctx context.Context, req *do.LessonLearnProgressReq) (*model.LessonLearnProgress, error) {
	qs := query.Use(r.db).LessonLearnProgress
	return qs.WithContext(ctx).Where(
		qs.UserID.Eq(req.UserID),
		qs.CourseID.Eq(req.CourseID),
		qs.LessonID.Eq(req.LessonID),
	).First()
}

func (r *LessonRepo) GetUserLearnProgresses(ctx context.Context, userID int64, offset, limit int) ([]*model.LessonLearnProgress, int64, error) {
	qs := query.Use(r.db).LessonLearnProgress
	tx := qs.WithContext(ctx).Where(qs.UserID.Eq(userID), qs.LearnStatus.Eq(consts.LearnStatusLearning))
	total, err := tx.Count()
	if err != nil {
		return nil, 0, err
	}
	list, err := tx.Order(qs.UpdateAt.Desc()).Offset(offset).Limit(limit).Find()
	return list, total, err
}

func (r *LessonRepo) UpsertLessonLearnProgress(ctx context.Context, req *do.LessonLearnProgressUpdate) error {
	now := time.Now()
	qs := query.Use(r.db).LessonLearnProgress
	progress := model.LessonLearnProgress{
		CourseID:     req.CourseID,
		LessonID:     req.LessonID,
		UserID:       req.UserID,
		PlayPosition: req.PlayPosition,
		LearnStatus:  consts.LearnStatusLearning,
		CreateAt:     now,
		UpdateAt:     now,
	}
	return qs.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: qs.UserID.ColumnName().String()},
			{Name: qs.CourseID.ColumnName().String()},
			{Name: qs.LessonID.ColumnName().String()},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			qs.PlayPosition.ColumnName().String(): req.PlayPosition,
			qs.LearnStatus.ColumnName().String():  consts.LearnStatusLearning,
			qs.UpdateAt.ColumnName().String():     now,
		}),
	}).Create(&progress)
}

func (r *LessonRepo) CreateLessonLearnRecord(ctx context.Context, req *do.LessonLearnRecordCreate) error {
	qs := query.Use(r.db).LessonLearnRecord
	return qs.WithContext(ctx).Create(&model.LessonLearnRecord{
		UserID:    req.UserID,
		CourseID:  req.CourseID,
		LessonID:  req.LessonID,
		EntryTime: req.EntryTime,
		ExitTime:  req.ExitTime,
		Duration:  int64(req.Duration),
		LastType:  req.LastType,
	})
}
