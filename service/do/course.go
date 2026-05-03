package do

import (
	"github.com/JunLang-7/mall/common"
)

type CreateCourse struct {
	UserID         int64
	Name           string
	CoursePrice    int64
	ServiceTime    int32
	LearnTime      int32
	Sort           int32
	Features       []string
	UpdateStatus   int32
	CoverKey       string
	DetailCoverKey string
	Detail         string
}

type UpdateCourse struct {
	UserID         int64
	ID             int64
	Name           string
	CoursePrice    int64
	ServiceTime    int32
	LearnTime      int32
	Sort           int32
	Features       []string
	UpdateStatus   int32
	CoverKey       string
	DetailCoverKey string
	Detail         string
}

type UpdateCourseStatus struct {
	UserID int64
	ID     int64
	Status int32
}

type CourseList struct {
	common.Pager
	ID              int64
	NameKW          string
	CreateStartTime int64
	CreateEndTime   int64
	UpdateStartTime int64
	UpdateEndTime   int64
	UpdateStatus    int32
	Status          int32
}

type AddCatalog struct {
	CourseID int64
	Name     string
	ParentID int64
	Sort     int32
	Level    int32
	UserID   int64
}

type UpdateCatalog struct {
	ID     int64
	UserID int64
	Name   string
}

type CatalogSort struct {
	ID       int64
	Sort     int32
	ParentID int64
	Level    int32
	Lessons  []*common.IDSort
}

type UpdateCatalogSort struct {
	SortList []*CatalogSort
	UserID   int64
}

type AddCatalogLesson struct {
	UserID    int64
	CourseID  int64
	CatalogID int64
	LessonIDs []int64
}

type UpdateCatalogLesson struct {
	ID          int64
	UserID      int64
	Name        string
	EnableTrial int32
	ShowTime    int64
}

type RemoveCatalogLesson struct {
	UserID int64
	IDs    []int64
}
